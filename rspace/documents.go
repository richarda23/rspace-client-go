package rspace

import (
	"encoding/json"
	"fmt"

	//	"net/url"
	//	"strconv"
	"time"
)

type DocumentService struct {
	BaseService
}

func (ds *DocumentService) documentsUrl() string {
	return ds.BaseUrl.String() + "/documents"
}

func (ds *DocumentService) statusUrl() string {
	return ds.BaseUrl.String() + "/status"
}

// GetStatus returns the result of the /status endpoint
func (ds *DocumentService) GetStatus() (*Status, error) {
	time.Sleep(ds.Delay)
	status, err := ds.doGet(ds.statusUrl())
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
	url := ds._generateUrl(config, "", false)
	return ds._doDocList(url)
}

func (ds *DocumentService) _doDocList(url string) (*DocumentList, error) {
	data, err := ds.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = DocumentList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
func (ds *DocumentService) _generateUrl(config RecordListingConfig, searchTerm string, isAdvancedSearch bool) string {
	params := config.toParams()
	if len(searchTerm) > 0 {
		if isAdvancedSearch {
			params.Add("advancedQuery", searchTerm)
		} else {
			params.Add("query", searchTerm)
		}
	}
	encoded := params.Encode()
	url := ds.documentsUrl() + "?" + encoded
	return url
}

//SearchDocuments performs basic search of a single search term, performing a global search
func (ds *DocumentService) SearchDocuments(config RecordListingConfig, searchTerm string) (*DocumentList, error) {
	time.Sleep(ds.Delay)
	url := ds._generateUrl(config, searchTerm, false)
	return ds._doDocList(url)
}

func (ds *DocumentService) AdvancedSearchDocuments(config RecordListingConfig, searchQuery *SearchQuery) (*DocumentList, error) {
	time.Sleep(ds.Delay)
	urlStr := ""
	if searchQuery != nil {
		queryJson, _ := json.Marshal(searchQuery)
		urlStr = ds._generateUrl(config, string(queryJson), true)
	} else {
		urlStr = ds._generateUrl(config, "", false)
	}
	return ds._doDocList(urlStr)
}

// DocumentById retrieves full document content
func (ds *DocumentService) DocumentById(docId int) (*Document, error) {
	time.Sleep(ds.Delay)
	url := fmt.Sprintf("%s/%d", ds.documentsUrl(), docId)
	data, err := ds.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = Document{}
	json.Unmarshal(data, &result)
	return &result, nil
}

// DeleteDocument attempts to delete the document with the specified ID
func (ds *DocumentService) DeleteDocument(documentId int) (bool, error) {
	time.Sleep(ds.Delay)
	url := fmt.Sprintf("%s/%d", ds.documentsUrl(), documentId)
	return ds.doDelete(url)
}

// DocumentNew creates a new RSpace document
func (ds *DocumentService) DocumentNew(post *DocumentPost) (*DocumentInfo, error) {
	time.Sleep(ds.Delay)
	data, err := ds.doPostJsonBody(post, ds.documentsUrl())
	if err != nil {
		return nil, err
	}
	result := &DocumentInfo{}
	json.Unmarshal(data, result)
	return result, nil
}

// NewBasicDocumentWithContent creates a new BasicDocument document with name, tags(optional) and content in a
// single text field.
func (ds *DocumentService) NewBasicDocumentWithContent(name string, tags string, contentHtml string) (*DocumentInfo, error) {
	time.Sleep(ds.Delay)
	post := BasicPost(name, tags)
	content := FieldContent{contentHtml}
	fields := make([]FieldContent, 1)
	fields[0] = content
	post.Fields = fields
	return ds.doPostCreateDocument(post)
}

// NewEmptyBasicDocument creates a new, empty BasicDocument with no content.
func (ds *DocumentService) NewEmptyBasicDocument(name string, tags string) (*DocumentInfo, error) {
	time.Sleep(ds.Delay)
	post := BasicPost(name, tags)
	return ds.doPostCreateDocument(post)
}

func (ds *DocumentService) doPostCreateDocument(postData *DocumentPost) (*DocumentInfo, error) {
	data, err := ds.doPostJsonBody(postData, ds.documentsUrl())
	if err != nil {
		return nil, err
	}
	result := &DocumentInfo{}
	json.Unmarshal(data, result)
	return result, nil
}
