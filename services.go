package rspace

import (
	"time"
	"net/url"
)
const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)
type BaseService struct {
	Delay time.Duration
	ApiKey string
	BaseUrl *url.URL

}
func baseService() BaseService {
	return BaseService{
		Delay: time.Duration(100) * time.Millisecond}
}
type RsWebClient struct {
	activityS *ActivityService
	documentS *DocumentService
	folderS *FolderService
	formS *FormService
	fileS *FileService
	sysS *SysadminService
}
var (
  webClient *RsWebClient 
)
func init (){
	url, _ :=  url.Parse(getenv(BASE_URL_ENV_NAME))
	webClient = NewWebClient(url, getenv(APIKEY_ENV_NAME)) 
}

func (wc *RsWebClient) activityService() *ActivityService {
	return wc.activityS
}
func (wc *RsWebClient) documentService() *DocumentService {
	return wc.documentS
}
func (wc *RsWebClient) folderService() *FolderService {
	return wc.folderS
}
func (wc *RsWebClient) fileService() *FileService {
	return wc.fileS
}
func (wc *RsWebClient) formService() *FormService {
	return wc.formS
}	
func (wc *RsWebClient) sysadminService() *SysadminService {
	return wc.sysS
}	


func NewWebClient (baseUrl *url.URL, apiKey string) *RsWebClient {
	base := baseService()
	base.ApiKey = apiKey
	base.BaseUrl= baseUrl
	wc := RsWebClient{}
	wc.activityS = &ActivityService {BaseService: base}
	wc.documentS = &DocumentService {BaseService: base}
	wc.folderS = &FolderService {BaseService: base}
	wc.formS = &FormService {BaseService: base}
	wc.fileS = &FileService {BaseService: base}
	return &wc
}

