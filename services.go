package rspace

import (
	"net/url"
	"time"
	"net/http"
	"io/ioutil"
	"io"
	"os"
	"encoding/json"
	"bytes"
)

var (
	webClient *RsWebClient
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

type BaseService struct {
	Delay   time.Duration
	ApiKey  string
	BaseUrl *url.URL
}

func (bs *BaseService) doPostJsonBody(post interface{}, urlString string) ([]byte, error) {
	time.Sleep(time.Duration(100) * time.Millisecond)
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", urlString, bytes.NewBuffer(formData))
	bs.addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err2 := testResponseForError(data, resp); err2 != nil {
		return nil, err2
	}
	return data, nil
}

// doDelete  attempts to delete a resource specified by the URL. If successful, returns true, else returns false, with a possible
// non-null error
func (bs *BaseService) doDelete(url string) (bool, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	bs.addAuthHeader(req)
	resp, e := client.Do(req)
	if e != nil {
		Log.Error(e)
		return false, e
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := testResponseForError(data, resp); err != nil {
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
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	bs.addAuthHeader(req)
	resp, err := client.Do(req)
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

// doGet makes an authenticated API request to a URL expecting a string
// response (typically JSON)
func (bs *BaseService) doGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	bs.addAuthHeader(req)
	resp, e := client.Do(req)
	if e != nil {
		Log.Error(e)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("resp:" + string(data))
	if err := testResponseForError(data, resp); err != nil {
		return nil, err
	}
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
}

func (ws *RsWebClient) Users(lastLoginBefore time.Time, creationDateBefore time.Time) (*UserList, error) {
	return ws.sysadminS.Users(lastLoginBefore, creationDateBefore)
}

// UserNew creates a new user account. Requires sysadmin role
func (ws *RsWebClient) UserNew (userPost *UserPost) (*UserInfo, error) {
	return ws.sysadminS.UserNew(userPost)
}

// GroupNew creates a new group with the specified users and PI. Requires sysadmin role
func (ws *RsWebClient) GroupNew (groupPost *GroupPost) (*GroupInfo, error) {
	return ws.sysadminS.GroupNew(groupPost)
}

// Forms returns a paginated listing of Forms
func (fs *RsWebClient) Forms(config RecordListingConfig) (*FormList, error) {
	return fs.formS.Forms(config)
}

//  Documents returns a paginated listing of RSpace documents
func (ds *RsWebClient) Documents (config RecordListingConfig) (*DocumentList, error){
	return ds.documentS.Documents(config)
}

// Activities queries the audit trail and returns a list of events.
func (ds *RsWebClient) Activities (query *ActivityQuery, pgCrit RecordListingConfig) (*ActivityList, error){
	return ds.activityS.Activities(query, pgCrit)
}

// SearchDocuments performs a global search for 'searchTerm' across all  searchable fields
func (ds *RsWebClient) SearchDocuments (config RecordListingConfig, searchTerm string) (*DocumentList, error){
	return ds.documentS.SearchDocuments(config, searchTerm)
}

// AdvancedSearchDocuments performs a search for the terms specified in 'searchQuery'
func (ds *RsWebClient) AdvancedSearchDocuments (config RecordListingConfig, searchQuery *SearchQuery) (*DocumentList, error){
	return ds.documentS.AdvancedSearchDocuments(config, searchQuery)
}

// Status returns simple information about the current server
func (ds *RsWebClient) Status () (*Status, error){
	return ds.documentS.GetStatus()
}

// NewEmptyBasicDocument creates a Basic (single text field) document with no content
func (ds *RsWebClient) NewEmptyBasicDocument (name, tags string) *DocumentInfo{
	return ds.documentS.NewEmptyBasicDocument(name, tags)
}
func (ds *RsWebClient) NewBasicDocumentWithContent (name, tags, content string) *DocumentInfo{
	return ds.documentS.NewBasicDocumentWithContent(name, tags, content)
}
// FolderTree returns a list of items in the specified folder
func (fs *RsWebClient) FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) (*FolderList, error) {
	return fs.folderS.FolderTree(config , folderId , typesToInclude ) 
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
	return &wc
}
