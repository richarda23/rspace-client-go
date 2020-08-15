package rspace

import (
	"encoding/json"
	"errors"
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

func (ex *ExportService) ExportSubmit(post ExportPost) (*Job, error) {
	url := ex.exportUrl()
	url = ex.makeUrl(post, url)
	var emptyBody struct{}
	data, err := ex.doPostJsonBody(emptyBody, url)
	if err != nil {
		return nil, err
	}
	var job = &Job{}
	json.Unmarshal(data, job)
	return job, nil
}

// Export does an export, blocking till job has finished.
// The returned job, if completed successfully, will contain a download link.
func (fs *ExportService) Export(post ExportPost) (*Job, error) {
	job, err := fs.ExportSubmit(post)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	initialSleepDuration := 100
	time.Sleep(time.Duration(initialSleepDuration) * time.Millisecond)
	for {
		job, err = fs.GetJob(job.Id)
		if err != nil {
			return nil, err
		}
		pc := job.PercentComplete
		if !job.IsTerminated() {
			if pc > 0 && pc < 100 {
				sleepMs, _ := calculateSleepTime(pc, start)
				time.Sleep(*sleepMs)
			}
		} else if job.IsCompleted() {
			Log.Infof("Completed, download link is %s", job.DownloadLink().String())
			break
		} else if job.IsTerminated() {
			Log.Infof("Job terminated unsuccessfully with status %s", job.Status)
			break
		}
	}
	return job, nil
}

// sleeps maximum of 3 seconds, or 1/5th of expected remaining time
func calculateSleepTime(pcComplete float32, start time.Time) (*time.Duration, error) {
	if pcComplete == 0 {
		return nil, errors.New("pcComplete must be > 0 to calculate sleep period")
	}
	elapsedTimeMs := float32(time.Now().Sub(start).Milliseconds())
	Log.Infof("elapsed time is %3.2f ms, pc = %.2f\n", elapsedTimeMs, pcComplete)
	expectedCompletionTime :=
		(elapsedTimeMs / pcComplete) * 100
	Log.Infof("expected completion time is %3.3f ms\n", expectedCompletionTime)

	sleepDurationF := math.Max(3000, float64(expectedCompletionTime-elapsedTimeMs)/5)
	Log.Infof("will sleep for %3.2f ms\n", sleepDurationF)
	duration := time.Duration(int64(sleepDurationF)) * time.Millisecond
	return &duration, nil
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
