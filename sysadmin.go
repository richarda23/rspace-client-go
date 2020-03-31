package rspace

import (
	"bytes"
	"encoding/json"
	"time"
	"net/http"
)

type SysadminService struct {
	BaseService
}

func systemUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/sysadmin"
}

// DocumentNew creates a new RSpace document
func (ds *SysadminService) UserNew(post *UserPost) *UserInfo {
	time.Sleep(ds.Delay)
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST",systemUrl()+"/users", bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		Log.Error(err)
	}
	result := &UserInfo{}
	Unmarshal(resp, result)
	return result
}

func (ds *SysadminService) GroupNew(post *GroupPost) (*GroupInfo, error) {
 //TODO
 return nil, nil
}

