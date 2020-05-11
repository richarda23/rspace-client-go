package rspace

import (
	"encoding/json"
	"strconv"
	"time"
)

type SharingService struct {
	BaseService
}

func (fs *SharingService) sharingUrl() string {
	return fs.BaseUrl.String() + "/share"
}

// Unshare unshares an item from a user or group.
// The id to be passed is the id of a ShareInfo, not the the Id of an RSpace document.
func (fs *SharingService) Unshare(shareId int) (bool, error) {
	time.Sleep(fs.Delay)
	resp, err := fs.doDelete(fs.sharingUrl() + "/" + strconv.Itoa(shareId))
	if resp == false {
		return false, err
	} else {
		return true, nil
	}
}

// Share an item with a group or user Id
func (fs *SharingService) Share(post *SharePost) (*ShareInfoList, error) {
	time.Sleep(fs.Delay)
	data, err := fs.doPostJsonBody(post, fs.sharingUrl())
	if err != nil {
		return nil, err
	}
	var result = ShareInfoList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
