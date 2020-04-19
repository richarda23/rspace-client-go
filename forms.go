package rspace

import (
	"encoding/json"
	"time"
)

type FormService struct {
	BaseService
}

func formsUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/forms"
}

// FormTree produces paginated listing of items in form
func (fs *FormService) Forms(config RecordListingConfig) (*FormList, error) {
	time.Sleep(fs.Delay)
	url := formsUrl()
	if paramStr := config.toParams().Encode(); len(paramStr) > 0 {
		url = url + "?" + paramStr
	}
	data, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = FormList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
