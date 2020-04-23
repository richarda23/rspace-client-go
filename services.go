package rspace

import (
	"net/url"
	"time"
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
