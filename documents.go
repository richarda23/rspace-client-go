package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type DocumentService struct {
	BaseService
}

func documentsUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/documents"
}

// GetStatus returns the result of the /status endpoint
func (ds *DocumentService) GetStatus() (*Status, error) {
	time.Sleep(ds.Delay)
	status, err := DoGet(getenv(BASE_URL_ENV_NAME) + "/status")
	if err != nil {
		return nil, err
	}
	res := Status{}
	json.Unmarshal(status, &res)
	return &res, nil
}

// Paginated listing of Documents
func (ds *DocumentService) Documents(config RecordListingConfig) (*DocumentList, error) {
	time.Sleep(ds.Delay)
	url := ds._generateUrl(config)
	return ds._doDocList(url)
}

func (ds *DocumentService) _doDocList(url string) (*DocumentList, error) {
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = DocumentList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
func (ds *DocumentService) _generateUrl(config RecordListingConfig) string {
	params := url.Values{}
	params.Add("pageSize", strconv.Itoa(config.PageSize))
	params.Add("pageNumber", strconv.Itoa(config.PageNumber))
	encoded := params.Encode()
	url := documentsUrl() + "?" + encoded
	return url
}

//SearchDocuments performs basic search of a single search term, performing a global search
func (ds *DocumentService) SearchDocuments(config RecordListingConfig, searchTerm string) (*DocumentList, error) {
	time.Sleep(ds.Delay)
	url := ds._generateUrl(config)
	if len(searchTerm) > 0 {
		url = url + "&query=" + searchTerm
	}
	return ds._doDocList(url)
}

func (ds *DocumentService) AdvancedSearchDocuments(config RecordListingConfig, searchQuery *SearchQuery) (*DocumentList, error) {
	time.Sleep(ds.Delay)
	urlStr := ds._generateUrl(config)

	if searchQuery != nil {
		queryJson, _ := json.Marshal(searchQuery)
		params := url.Values{}
		params.Add("advancedQuery", string(queryJson))
		encoded := params.Encode()
		urlStr = urlStr + "&" + encoded
	}
	return ds._doDocList(urlStr)
}

// DocumentById retrieves full document content
func (ds *DocumentService) DocumentById(docId int) (*Document, error) {
	time.Sleep(ds.Delay)
	url := fmt.Sprintf("%s/%d", documentsUrl(), docId)
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = Document{}
	json.Unmarshal(data, &result)
	return &result, nil
}

// DocumentNew creates a new RSpace document
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
	time.Sleep(ds.Delay)
	post := BasicPost(name, tags)
	content := FieldContent{contentHtml}
	fields := make([]FieldContent, 1)
	fields[0] = content
	post.Fields = fields
	return doPostCreateDocument(post)
}

// NewEmptyBasicDocument creates a new, empty BasicDocument with no content.
func (ds *DocumentService) NewEmptyBasicDocument(name string, tags string) *DocumentInfo {
	time.Sleep(ds.Delay)
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
