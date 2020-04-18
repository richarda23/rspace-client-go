package rspace

import (
	"fmt"
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
	ActivityS *ActivityService
	documentS *DocumentService
	folderS   *FolderService
	FormS     *FormService
	fileS     *FileService
	SysadminS *SysadminService
}

func (ds *RsWebClient) Documents (config RecordListingConfig) (*DocumentList, error){
	return ds.documentS.Documents(config)
}

func (ds *RsWebClient) SearchDocuments (config RecordListingConfig, searchTerm string) (*DocumentList, error){
	return ds.documentS.SearchDocuments(config, searchTerm)
}

func (ds *RsWebClient) AdvancedSearchDocuments (config RecordListingConfig, searchQuery *SearchQuery) (*DocumentList, error){
	return ds.documentS.AdvancedSearchDocuments(config, searchQuery)
}

func (ds *RsWebClient) Status () (*Status, error){
	return ds.documentS.GetStatus()
}

func (ds *RsWebClient) NewEmptyBasicDocument (name string, tags string) *DocumentInfo{
	return ds.documentS.NewEmptyBasicDocument(name, tags)
}

func (fs *RsWebClient) FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) (*FolderList, error) {
	return fs.folderS.FolderTree(config , folderId , typesToInclude ) 
}
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
func (fs *RsWebClient) DeleteFolder(folderId int) (bool, error) {
	return fs.folderS.DeleteFolder(folderId)
}
func (fs *RsWebClient) FolderNew(post *FolderPost) (*Folder, error) {
	return fs.folderS.FolderNew(post)
}

func NewWebClient(baseUrl *url.URL, apiKey string) *RsWebClient {
	fmt.Println("In newwebcient")
	base := baseService()
	base.ApiKey = apiKey
	base.BaseUrl = baseUrl
	wc := RsWebClient{}
	wc.ActivityS = &ActivityService{BaseService: base}
	wc.documentS = &DocumentService{BaseService: base}
	wc.folderS = &FolderService{BaseService: base}
	wc.FormS = &FormService{BaseService: base}
	wc.fileS = &FileService{BaseService: base}
	wc.SysadminS = &SysadminService{BaseService: base}
	return &wc
}
