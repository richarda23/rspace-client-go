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
	DocumentS *DocumentService
	folderS   *FolderService
	FormS     *FormService
	FileS     *FileService
	SysadminS *SysadminService
}
func (fs *RsWebClient) FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) (*FolderList, error) {
	return fs.folderS.FolderTree(config , folderId , typesToInclude ) 
}
func (fs *RsWebClient) FolderById(folderId int) (*Folder, error) {
	return fs.folderS.FolderById(folderId)
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
	wc.DocumentS = &DocumentService{BaseService: base}
	wc.folderS = &FolderService{BaseService: base}
	wc.FormS = &FormService{BaseService: base}
	wc.FileS = &FileService{BaseService: base}
	return &wc
}
