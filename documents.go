package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func documentsUrl() string {
    return getenv(BASE_URL_ENV_NAME) + "/documents"
}
//GetStatus returns the result of the /status endpoint
func GetStatus() *Status {
	statusStr := DoGet(getenv(BASE_URL_ENV_NAME) + "/status")
	fmt.Println(statusStr)
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	return &res
}

// Paginated listing of Documents
func Documents(config RecordListingConfig) *DocumentList {
	url := documentsUrl() + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson := DoGet(url)
	var result = DocumentList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

//DocumentById retrieves full document content
func DocumentById(docId int) *Document {
	url := fmt.Sprintf("%s/%d", documentsUrl(), docId)
	docJson := DoGet(url)
	fmt.Println(docJson)
	var result = Document{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

//DocumentNew creates a new RSpace document
func DocumentNew(post *DocumentPost) *DocumentInfo {
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", documentsUrl(), bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}

// NewBasicDocumentWithContent creates a new BasicDocument document with name, tags(optional) and content in a 
// single text field.
func NewBasicDocumentWithContent(name string, tags string, contentHtml string) *DocumentInfo {
	post := BasicPost(name, tags)
	content := FieldContent{contentHtml}
	fields := make([]FieldContent, 1)
	fields[0] = content
	post.Fields = fields
	return doPostCreateDocument(post)
}

// NewEmptyBasicDocument creates a new, empty BasicDocument with no content.
func NewEmptyBasicDocument(name string, tags string) *DocumentInfo {
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
		log.Fatalln(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}

