package rspace

import (
	"encoding/json"
)

type FormService struct {
	BaseService
}

func (fs *FormService) formsUrl() string {
	return fs.BaseUrl.String() + "/forms"
}

// FormTree produces paginated listing of items in form
func (fs *FormService) Forms(config RecordListingConfig) (*FormList, error) {
	url := fs.formsUrl()
	if paramStr := config.toParams().Encode(); len(paramStr) > 0 {
		url = url + "?" + paramStr
	}
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = FormList{}
	json.Unmarshal(data, &result)
	return &result, nil
}
