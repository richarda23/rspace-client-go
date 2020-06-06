package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	webClient *RsWebClient
)

const (
	APIKEY_ENV_NAME      = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME    = "RSPACE_URL"
	RATE_LIMIT_WAIT_TIME = "X-Rate-Limit-WaitTimeMillis"
)

type BaseService struct {
	Delay   time.Duration
	ApiKey  string
	BaseUrl *url.URL
}

func (bs *BaseService) doPutJsonBody(post interface{}, urlString string) ([]byte, error) {
	return bs._postOrPutJsonBody(post, urlString, "PUT")
}

func (bs *BaseService) doPostJsonBody(post interface{}, urlString string) ([]byte, error) {
	return bs._postOrPutJsonBody(post, urlString, "POST")
}

func (bs *BaseService) _postOrPutJsonBody(post interface{}, urlString, httpVerb string) ([]byte, error) {
	formData, _ := json.Marshal(post)
	//	Log.Info(string(formData))
	hc := http.Client{Timeout: time.Duration(15) * time.Second}
	req, err := http.NewRequest(httpVerb, urlString, bytes.NewBuffer(formData))
	bs.addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	retry := NewResilientClient(&hc)
	resp, err := retry.Do(req)
	if err != nil {
		return nil, err
	}
	if err2 := testResponseForError(resp); err2 != nil {
		return nil, err2
	}
	data, _ := ioutil.ReadAll(resp.Body)
	//Log.Debug(string(data))
	return data, nil
}

// doDelete  attempts to delete a resource specified by the URL. If successful, returns true, else returns false, with a possible
// non-null error
func (bs *BaseService) doDelete(url string) (bool, error) {
	client := HttpClientNew(10)
	retry := NewResilientClient(client)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	bs.addAuthHeader(req)
	resp, e := retry.Do(req)
	if e != nil {
		Log.Error(e)
		return false, e
	}
	if err := testResponseForError(resp); err != nil {
		return false, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	} else {
		return false, nil
	}
}

// DoGetToFile saves the response from an HTTP GET request to the specified file.
// If the response fails or the file cannot be created returns an error.
// 'filepath' argument should be absolute path to a file. If the file exists, it will be overwritten. If it doesn't exist, it will be created.
func (bs *BaseService) doGetToFile(url string, filepath string) error {
	client := HttpClientNew(10)
	retry := NewResilientClient(client)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	bs.addAuthHeader(req)
	resp, err := retry.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	return err
}

// RateLimitData stores information received in HTTP Response Headers
// about API usage rates. If this information could not be retrieved from a response then
// value will be -100
type RateLimitData struct {
	WaitTimeMillis int
}

// Stringer implementation
func (rld RateLimitData) String() string {
	return fmt.Sprintf("Wait time till next request: %d",
		rld.WaitTimeMillis)
}

func NewRateLimitData(resp *http.Response) RateLimitData {
	errorValue := -100
	rl, err := strconv.Atoi(resp.Header.Get(RATE_LIMIT_WAIT_TIME))
	if err != nil {
		rl = errorValue
	}
	return RateLimitData{rl}
}

// doGet makes an authenticated API request to a URL expecting a string
// response (typically JSON)
func (bs *BaseService) doGet(url string) ([]byte, error) {
	client := HttpClientNew(15)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	bs.addAuthHeader(req)
	retry := NewResilientClient(client)
	resp, e := retry.Do(req)
	if e != nil {
		Log.Error(e)
		return nil, e
	}
	//	var rld RateLimitData = NewRateLimitData(resp)
	//	Log.Info(rld.String())

	//fmt.Println("resp:" + string(data))
	if err := testResponseForError(resp); err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	return data, nil
}

func (bs *BaseService) addAuthHeader(req *http.Request) {
	req.Header.Add("apiKey", bs.ApiKey)
}

func baseService() BaseService {
	return BaseService{
		Delay: time.Duration(100) * time.Millisecond}
}

type RsWebClient struct {
	activityS *ActivityService
	documentS *DocumentService
	folderS   *FolderService
	formS     *FormService
	fileS     *FileService
	sysadminS *SysadminService
	importS   *ImportService
	groupS    *GroupService
	sharingS  *SharingService
	exportS   *ExportService
}

func (ws *RsWebClient) Groups() (*GroupList, error) {
	return ws.groupS.Groups()
}
func (ws *RsWebClient) Users(lastLoginBefore time.Time, creationDateBefore time.Time, cfg RecordListingConfig) (*UserList, error) {
	return ws.sysadminS.Users(lastLoginBefore, creationDateBefore, cfg)
}

// UserNew creates a new user account. Requires sysadmin role
func (ws *RsWebClient) UserNew(userPost *UserPost) (*UserInfo, error) {
	return ws.sysadminS.UserNew(userPost)
}

// GroupNew creates a new group with the specified users and PI. Requires sysadmin role
func (ws *RsWebClient) GroupNew(groupPost *GroupPost) (*GroupInfo, error) {
	return ws.sysadminS.GroupNew(groupPost)
}

// Forms returns a paginated listing of Forms
func (fs *RsWebClient) Forms(config RecordListingConfig) (*FormList, error) {
	return fs.formS.Forms(config, "")
}

// FormSearch returns a paginated listing of Forms filtered by optional search query
func (fs *RsWebClient) FormSearch(config RecordListingConfig, query string) (*FormList, error) {
	return fs.formS.Forms(config, query)
}

//  Documents returns a paginated listing of RSpace documents
func (ds *RsWebClient) Documents(config RecordListingConfig) (*DocumentList, error) {
	return ds.documentS.Documents(config)
}

// Activities queries the audit trail and returns a list of events.
func (ds *RsWebClient) Activities(query *ActivityQuery, pgCrit RecordListingConfig) (*ActivityList, error) {
	return ds.activityS.Activities(query, pgCrit)
}

// SearchDocuments performs a global search for 'searchTerm' across all  searchable fields
func (ds *RsWebClient) SearchDocuments(config RecordListingConfig, searchTerm string) (*DocumentList, error) {
	return ds.documentS.SearchDocuments(config, searchTerm)
}

// AdvancedSearchDocuments performs a search for the terms specified in 'searchQuery'
func (ds *RsWebClient) AdvancedSearchDocuments(config RecordListingConfig, searchQuery *SearchQuery) (*DocumentList, error) {
	return ds.documentS.AdvancedSearchDocuments(config, searchQuery)
}

// Status returns simple information about the current server
func (ds *RsWebClient) Status() (*Status, error) {
	return ds.documentS.GetStatus()
}

// NewEmptyBasicDocument creates a Basic (single text field) document with no content
func (ds *RsWebClient) NewEmptyBasicDocument(name, tags string) (*Document, error) {
	return ds.documentS.NewEmptyBasicDocument(name, tags)
}
func (ds *RsWebClient) NewBasicDocumentWithContent(name, tags, content string) (*Document, error) {
	return ds.documentS.NewBasicDocumentWithContent(name, tags, content)
}

func (ds *RsWebClient) NewDocumentWithContent(docPost *DocumentPost) (*Document, error) {
	return ds.documentS.DocumentNew(docPost)
}

func (ds *RsWebClient) DocumentEdit(docId int, docPost *DocumentPost) (*Document, error) {
	return ds.documentS.DocumentEdit(docId, docPost)
}

func (ds *RsWebClient) DocumentById(docId int) (*Document, error) {
	return ds.documentS.DocumentById(docId)
}

// FolderTree returns a list of items in the specified folder
func (fs *RsWebClient) FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) (*FolderList, error) {
	return fs.folderS.FolderTree(config, folderId, typesToInclude)
}

// FolderById returns information about the folder specified by folderId
func (fs *RsWebClient) FolderById(folderId int) (*Folder, error) {
	return fs.folderS.FolderById(folderId)
}

// Uploads a single file
func (fs *RsWebClient) UploadFile(path string) (*FileInfo, error) {
	return fs.fileS.UploadFile(path)
}

// Lists Gallery files, optionally filtered by a media type
func (fs *RsWebClient) Files(config RecordListingConfig, mediaType string) (*FileList, error) {
	return fs.fileS.Files(config, mediaType)
}

// Lists Gallery files, optionally filtered by a media type
func (fs *RsWebClient) FileById(id int) (*FileInfo, error) {
	return fs.fileS.FileById(id)
}

func (fs *RsWebClient) UploadFileNewVersion(path string, fileToReplaceId int) (*FileInfo, error) {
	return fs.fileS.UploadFileNewVersion(path, fileToReplaceId)
}

//Download downloads a file attachment with the given ID to the location set by the path.
func (fs *RsWebClient) Download(id int, path string) (*FileInfo, error) {
	return fs.fileS.DownloadFile(id, path)
}

// DeleteFolder deletes the given folder
func (fs *RsWebClient) DeleteFolder(folderId int) (bool, error) {
	return fs.folderS.DeleteFolder(folderId)
}

// FolderNew creates a new folder or notebook
func (fs *RsWebClient) FolderNew(post *FolderPost) (*Folder, error) {
	return fs.folderS.FolderNew(post)
}

func (fs *RsWebClient) ImportWord(path string, folderId int, imageFolderId int) (*DocumentInfo, error) {
	return fs.importS.ImportWord(path, folderId, imageFolderId)
}

// Share shares one or more items with one or more groups and users.
// Sharer and sharee must have a group in common.
func (client *RsWebClient) Share(post *SharePost) (*ShareInfoList, error) {
	return client.sharingS.Share(post)
}
func (client *RsWebClient) Unshare(shareId int) (bool, error) {
	return client.sharingS.Unshare(shareId)
}

func (client *RsWebClient) ShareList(query string, cfg RecordListingConfig) (*SharedItemList, error) {
	return client.sharingS.SharedItemList(query, cfg)
}

func (client *RsWebClient) Export(post ExportPost) (*Job, error) {
	return client.exportS.Export(post)
}

func NewWebClient(baseUrl *url.URL, apiKey string) *RsWebClient {
	base := baseService()
	base.ApiKey = apiKey
	base.BaseUrl = baseUrl
	wc := RsWebClient{}
	wc.activityS = &ActivityService{BaseService: base}
	wc.documentS = &DocumentService{BaseService: base}
	wc.folderS = &FolderService{BaseService: base}
	wc.formS = &FormService{BaseService: base}
	wc.fileS = &FileService{BaseService: base}
	wc.sysadminS = &SysadminService{BaseService: base}
	wc.importS = &ImportService{BaseService: base}
	wc.groupS = &GroupService{BaseService: base}
	wc.sharingS = &SharingService{BaseService: base}
	wc.exportS = &ExportService{BaseService: base}

	return &wc
}
