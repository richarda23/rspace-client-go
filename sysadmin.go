package rspace

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	data, err := doPostJsonBody(post, systemUrl()+"/users")
	if err != nil {
		return nil, err
	}
	result := &UserInfo{}
	json.Unmarshal(data, result)
	return result, nil
}

func (ds *SysadminService) Users(lastLoginBefore time.Time, creationDateBefore time.Time) (*UserList, error) {
	time.Sleep(ds.Delay)
	params := url.Values{}
	params.Add("tempAccountsOnly", strconv.FormatBool(false))
	if !lastLoginBefore.IsZero() {
		params.Add("lastLoginBefore", lastLoginBefore.Format("2006-02-01"))
	}
	if !creationDateBefore.IsZero() {
		params.Add("createdBefore", creationDateBefore.Format("2006-02-01"))
	}
	encoded := params.Encode()
	url := systemUrl() + "/users?" + encoded
	fmt.Println(url)
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	rc := &UserList{}
	json.Unmarshal(data, rc)
	return rc, err

}

func (ds *SysadminService) GroupNew(post *GroupPost) (*GroupInfo, error) {
	time.Sleep(ds.Delay)
	data, err := doPostJsonBody(post, systemUrl()+"/groups")
	if err != nil {
		return nil, err
	}
	result := &GroupInfo{}
	json.Unmarshal(data, result)
	return result, nil
}
