package rspace

import (
	"encoding/json"
	"fmt"
)

type ExportService struct {
	BaseService
}

func (fs *ExportService) exportUrl() string {
	return fs.BaseUrl.String() + "/export"
}

// Forms produces paginated listing of items in form
func (fs *ExportService) Export(post ExportPost) (*Job, error) {
	url := fs.exportUrl()
	url = fs.makeUrl(post, url)
	var q struct{}
	data, err := fs.doPostJsonBody(q, url)
	if err != nil {
		return nil, err
	}
	var result = Job{}
	json.Unmarshal(data, &result)
	return &result, nil
}

func (es *ExportService) makeUrl(post ExportPost, baseUrl string) string {
	url := baseUrl
	if post.Id == 0 {
		url = fmt.Sprintf("%s/%s/%s", url, post.Format.String(),
			post.Scope.String())
	} else {
		url = fmt.Sprintf("%s/%s/%s/%d", url, post.Format.String(),
			post.Scope.String(), post.Id)
	}
	return url
}
