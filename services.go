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
	FolderS   *FolderService
	FormS     *FormService
	FileS     *FileService
	SysadminS *SysadminService
}

func NewWebClient(baseUrl *url.URL, apiKey string) *RsWebClient {
	fmt.Println("In newwebcient")
	base := baseService()
	base.ApiKey = apiKey
	base.BaseUrl = baseUrl
	wc := RsWebClient{}
	wc.ActivityS = &ActivityService{BaseService: base}
	wc.DocumentS = &DocumentService{BaseService: base}
	wc.FolderS = &FolderService{BaseService: base}
	wc.FormS = &FormService{BaseService: base}
	wc.FileS = &FileService{BaseService: base}
	return &wc
}
