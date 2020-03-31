package rspace

import (
	"encoding/json"
	"time"
)

type SysadminService struct {
	BaseService
}

func systemUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/sysadmin"
}

// DocumentNew creates a new RSpace document
func (ds *SysadminService) UserNew(post *UserPost) (*UserInfo, error) {
	time.Sleep(ds.Delay)
	data, err := doPostJsonBody (post, systemUrl()+"/users")
	if err != nil {
		return nil, err
	}
	result := &UserInfo{}
	json.Unmarshal(data, result)
	return result,nil
}

func (ds *SysadminService) GroupNew(post *GroupPost) (*GroupInfo, error) {
	time.Sleep(ds.Delay)
	data, err := doPostJsonBody (post, systemUrl()+"/groups")
	if err != nil {
		return nil, err
	}
	result := &GroupInfo{}
	json.Unmarshal(data, result)
	return result,nil
}

