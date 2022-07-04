package rspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type FormService struct {
	BaseService
}

func (fs *FormService) formsUrl() string {
	return fs.BaseUrl.String() + "/forms"
}

func (fs *FormService) PublishForm(formId int) (*Form, error) {
	if formId <= 0 {
		return nil, errors.New("Form Id must be >= 1")
	}
	url := fmt.Sprintf("%s/%d/%s", fs.formsUrl(), formId, "publish")
	result, err := fs.doPut(url)
	if err != nil {
		return nil, err
	}
	var rc = &Form{}
	json.Unmarshal(result, &rc)
	return rc, nil

}

func (fs *FormService) CreateFormJson(jsonFormDef io.Reader) (*Form, error) {
	jsonBytes, _ := ioutil.ReadAll(jsonFormDef)
	result, err := fs.postOrPutJsonBodyBytes(jsonBytes, fs.formsUrl(), "POST")
	if err != nil {
		return nil, err
	}
	var rc = &Form{}
	json.Unmarshal(result, &rc)
	return rc, nil

}

func (fs *FormService) CreateFormYaml(yamlFormDef io.Reader) (*Form, error) {
	byteRes, err := ioutil.ReadAll(yamlFormDef)

	if err != nil {
		return nil, err
	}
	j2, err := yaml.YAMLToJSON(byteRes)
	if err != nil {
		return nil, err
	}
	return fs.CreateFormJson(bytes.NewBuffer(j2))

}

// Forms produces paginated listing of items in form
func (fs *FormService) Forms(config RecordListingConfig, query string) (*FormList, error) {
	url := fs.formsUrl()
	params := config.toParams()
	if len(query) > 0 {
		params.Set("query", query)
	}
	params.Set("scope", "all")

	if paramStr := params.Encode(); len(paramStr) > 0 {
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
