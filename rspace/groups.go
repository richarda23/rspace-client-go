package rspace

import (
	"encoding/json"
)

type GroupService struct {
	BaseService
}

func (fs *GroupService) groupsUrl() string {
	return fs.BaseUrl.String() + "/groups"
}

// FormTree produces paginated listing of items in form
func (fs *GroupService) Groups() (*GroupList, error) {
	url := fs.groupsUrl()
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = make([]*GroupInfo, 0)
	json.Unmarshal(data, &result)
	return &GroupList{Groups: result}, nil
}
