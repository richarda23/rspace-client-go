package rspace

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
	//	"fmt"
)

type SysadminService struct {
	BaseService
}

func (sysSvc *SysadminService) systemUrl() string {
	return sysSvc.BaseUrl.String() + "/sysadmin"
}

// UserNew creates a new user account.
func (ds *SysadminService) UserNew(post *UserPost) (*UserInfo, error) {
	data, err := ds.doPostJsonBody(post, ds.systemUrl()+"/users")
	if err != nil {
		return nil, err
	}
	result := &UserInfo{}
	json.Unmarshal(data, result)
	return result, nil
}

//Users lists users' biographical information
func (ds *SysadminService) Users(lastLoginBefore time.Time, creationDateBefore time.Time) (*UserList, error) {
	params := url.Values{}
	params.Add("tempAccountsOnly", strconv.FormatBool(false))
	if !lastLoginBefore.IsZero() {
		params.Add("lastLoginBefore", lastLoginBefore.Format("2006-01-02"))
	}
	if !creationDateBefore.IsZero() {
		params.Add("createdBefore", creationDateBefore.Format("2006-01-02"))
	}
	encoded := params.Encode()
	url := ds.BaseUrl.String() + "/sysadmin/users"
	if len(encoded) > 0 {
		url = url + "?" + encoded
	}
	data, err := ds.doGet(url)
	if err != nil {
		return nil, err
	}
	rc := &UserList{}
	json.Unmarshal(data, rc)
	return rc, err

}

//GroupNew creates a new group from existing users
func (ds *SysadminService) GroupNew(post *GroupPost) (*GroupInfo, error) {
	data, err := ds.doPostJsonBody(post, ds.systemUrl()+"/groups")
	if err != nil {
		return nil, err
	}
	result := &GroupInfo{}
	json.Unmarshal(data, result)
	return result, nil
}
