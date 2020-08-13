package rspace

import (
	"encoding/json"
	"fmt"
	"math"
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
func (fs *ExportService) Export(post ExportPost) (*Job, error) {
	url := fs.exportUrl()
	url = fs.makeUrl(post, url)
	var emptyBody struct{}
	start := time.Now()
	data, err := fs.doPostJsonBody(emptyBody, url)
	if err != nil {
		return nil, err
	}
	var job = &Job{}
	json.Unmarshal(data, job)
	initialSleepDuration := 100
	sleepDuration := initialSleepDuration
	time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
	for i := 0; i < 100; i++ {
		job, err = fs.GetJob(job.Id)
		if err != nil {
			return nil, err
		}
		pc := job.PercentComplete
		if !job.IsCompleted() {
			if pc > 0 && pc < 100 {
				elapsedTimeMs := float32(time.Now().Sub(start).Milliseconds())
				fmt.Printf("elapsed time is %3.2f ms, pc = %.2f\n", elapsedTimeMs, pc)
				expectedCompletionTime :=
					(elapsedTimeMs / pc) * 100
				fmt.Printf("expected completion time is %3.3f ms\n", expectedCompletionTime)

				sleepDurationF := math.Max(3000, float64(expectedCompletionTime-elapsedTimeMs)/5)
				fmt.Printf("will sleep for %3.2f ms\n", sleepDurationF)
				sleepDuration = int(sleepDurationF)
				time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
			}
		} else {
			Log.Infof("Completed, download link is %s", job.DownloadLink().String())
			break
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
