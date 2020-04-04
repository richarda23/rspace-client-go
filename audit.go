package rspace

import (
	"encoding/json"
	"fmt"
	"time"
)

type ActivityService struct {
	BaseService
}

func auditUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/activity"
}

// FolderTree produces paginated listing of items in folder. If folderId is 0 then Home Folder is lister
func (fs *ActivityService) Activities() (*ActivityList, error) {
	time.Sleep(fs.Delay)
	url := auditUrl()
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	var result = ActivityList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
