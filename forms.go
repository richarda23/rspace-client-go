package rspace

import (
	"encoding/json"
	"strconv"
	"time"
)

type FormService struct {
	BaseService
}

func formsUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/forms"
}

// FormTree produces paginated listing of items in form. If formId is 0 then Home Form is lister
func (fs *FormService) Forms(config RecordListingConfig) (*FormList, error) {
	time.Sleep(fs.Delay)
	url := formsUrl()
	url = url + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = FormList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
