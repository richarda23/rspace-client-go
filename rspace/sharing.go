package rspace

import (
	"encoding/json"
	"time"
)

type SharingService struct {
	BaseService
}

func (fs *SharingService) sharingUrl() string {
	return fs.BaseUrl.String() + "/share"
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
