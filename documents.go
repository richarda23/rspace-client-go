package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type DocumentService struct {
  BaseService
}

func documentsUrl() string {
    return getenv(BASE_URL_ENV_NAME) + "/documents"
}
//GetStatus returns the result of the /status endpoint
func (ds *DocumentService) GetStatus() (*Status, error) {
	time.Sleep(ds.Delay )
	statusStr, err := DoGet(getenv(BASE_URL_ENV_NAME) + "/status")
	if err != nil {
		return nil, err
	}
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	return &res, nil
}

// Paginated listing of Documents
func (ds *DocumentService) Documents(config RecordListingConfig) (*DocumentList, error) {
	time.Sleep(ds.Delay )
	url := documentsUrl() + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = DocumentList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result, nil
}

//DocumentById retrieves full document content
func (ds *DocumentService) DocumentById(docId int) (*Document, error)  {
	time.Sleep(ds.Delay)
	url := fmt.Sprintf("%s/%d", documentsUrl(), docId)
	docJson, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	fmt.Println(docJson)
	var result = Document{}
	json.Unmarshal([]byte(docJson), &result)
	return &result, nil
}

//DocumentNew creates a new RSpace document
func (ds *DocumentService) DocumentNew(post *DocumentPost) *DocumentInfo {
	time.Sleep(ds.Delay)
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", documentsUrl(), bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		Log.Error(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}

// NewBasicDocumentWithContent creates a new BasicDocument document with name, tags(optional) and content in a 
// single text field.
func (ds *DocumentService) NewBasicDocumentWithContent(name string, tags string, contentHtml string) *DocumentInfo {
	time.Sleep(ds.Delay )
	post := BasicPost(name, tags)
	content := FieldContent{contentHtml}
	fields := make([]FieldContent, 1)
	fields[0] = content
	post.Fields = fields
	return doPostCreateDocument(post)
}

// NewEmptyBasicDocument creates a new, empty BasicDocument with no content.
func (ds *DocumentService) NewEmptyBasicDocument(name string, tags string) *DocumentInfo {
	time.Sleep(ds.Delay )
	post := BasicPost(name, tags)
	return doPostCreateDocument(post)
}

func doPostCreateDocument(postData *DocumentPost) *DocumentInfo {
	hc := http.Client{}
	formData, _ := json.Marshal(postData)
	req, err := http.NewRequest("POST", documentsUrl(), bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		Log.Error(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}

