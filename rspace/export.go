package rspace

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExportService struct {
	BaseService
}

func (fs *ExportService) exportUrl() string {
	return fs.BaseUrl.String() + "/export"
}
func (fs *ExportService) jobsUrl() string {
	return fs.BaseUrl.String() + "/jobs"
}

func (fs *ExportService) GetJob(jobId int) (*Job, error) {
	url := fmt.Sprintf("%s/%d", fs.jobsUrl(), jobId)
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	jobrc := Job{}
	json.Unmarshal(data, &jobrc)
	return &jobrc, nil

}

// Export does an export, blocking till job has finished.
// The returned job will contain a download link.
// TODO fix 303 handling, add location header in RSpace
func (fs *ExportService) Export(post ExportPost) (*Job, error) {
	url := fs.exportUrl()
	url = fs.makeUrl(post, url)
	var emptyBody struct{}
	data, err := fs.doPostJsonBody(emptyBody, url)
	if err != nil {
		return nil, err
	}
	var job = &Job{}
	json.Unmarshal(data, job)
	for i := 0; i < 10; i++ {
		if !job.IsCompleted() {
			Log.Infof("Waiting for result %d", i)
			time.Sleep(time.Duration(3) * time.Second)
			job, err = fs.GetJob(job.Id)
			if err != nil {
				return nil, err
			}
		} else {
			Log.Infof("Completed, download link is %s", job.DownloadLink().String())
		}
	}

	return job, nil
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
