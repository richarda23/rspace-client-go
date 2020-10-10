package rspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (es *ExportService) GetJob(jobId int) (*Job, error) {
	url := fmt.Sprintf("%s/%d", es.jobsUrl(), jobId)
	data, err := es.doGet(url)
	if err != nil {
		return nil, err
	}
	jobrc := Job{}
	json.Unmarshal(data, &jobrc)
	return &jobrc, nil
}

//Download export downloads to the supplied filepath on local device
func (es *ExportService) DownloadExport(url string, outWriter io.Writer) error {
	return es.doGetToFile(url, outWriter)
}

func (ex *ExportService) exportSubmit(post ExportPost) (*Job, error) {
	url := ex.exportUrl()
	url, e := ex.makeUrl(post, url)
	if e != nil {
		return nil, e
	}

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
// reporter callback can be used to report or log progress  to client
func (fs *ExportService) Export(post ExportPost, waitForComplete bool, reporter func(string)) (*Job, error) {
	job, err := fs.exportSubmit(post)
	if !waitForComplete {
		return job, err
	}
	if err != nil {
		return nil, err
	}
	reporter(fmt.Sprintf("Waiting for job id %d: %s", job.Id, job.Links[0].Link))
	start := time.Now()
	initialSleepDuration := 1
	time.Sleep(time.Duration(initialSleepDuration) * time.Second)
	for {
		job, err = fs.GetJob(job.Id)
		if err != nil {
			return nil, err
		}
		pc := job.PercentComplete
		if !job.IsTerminated() {
			// if pc is 0, we have no info on which to base calculation
			if pc > 0 && pc < 100 {
				sleepMs, _ := calculateSleepTime(pc, start, reporter)
				time.Sleep(*sleepMs)
			} else {
				// a long job might still be on 0% progress for a while.
				time.Sleep(time.Duration(3 * time.Second))
			}
		} else if job.IsCompleted() {
			reporter(fmt.Sprintf("Completed, download link is %s", job.DownloadLink().String()))
			break
		} else if job.IsTerminated() {
			Log.Infof("Job terminated unsuccessfully with status %s", job.Status)
			break
		}
	}
	return job, nil
}

// sleeps maximum of 3 seconds, or 1/5th of expected remaining time
func calculateSleepTime(pcComplete float32, start time.Time, progressReporter func(string)) (*time.Duration, error) {
	if pcComplete == 0 {
		return nil, errors.New("pcComplete must be > 0 to calculate sleep period")
	}
	elapsedTimeS := float32(time.Now().Sub(start).Seconds())
	//progressReporter(fmt.Sprintf("elapsed time is %3.2f ms, pc = %.2f\n", elapsedTimeS, pcComplete)
	expectedCompletionTime :=
		((elapsedTimeS / pcComplete) * 100) - elapsedTimeS
	var minSleepTime float64 = 3.0
	sleepDurationF := math.Max(minSleepTime, float64(expectedCompletionTime-elapsedTimeS)/5)
	progressReporter(fmt.Sprintf("%3.1f%% complete. "+
		"Estimated remaining time: %3.1fs. Polling again in %4.1fs",
		pcComplete, expectedCompletionTime, sleepDurationF))
	duration := time.Duration(int64(sleepDurationF)) * time.Second
	return &duration, nil
}

func (es *ExportService) makeUrl(post ExportPost, baseUrl string) (string, error) {
	url := baseUrl
	if post.Id == 0 && !(post.Scope == SELECTION_EXPORT_SCOPE) {
		url = fmt.Sprintf("%s/%s/%s", url, post.Format.String(),
			post.Scope.String())
	} else if post.Scope == SELECTION_EXPORT_SCOPE {
		if len(post.ItemIds) == 0 {
			return "", errors.New("For selection scope, must supply >= 1 id")
		}
		url = fmt.Sprintf("%s/%s/%s?selections=%s&maxLinkLevel=%d", url, post.Format.String(),
			post.Scope.String(), post.ItemIdsToRequest(), post.MaxLinkLevel)
	} else {
		url = fmt.Sprintf("%s/%s/%s/%d", url, post.Format.String(),
			post.Scope.String(), post.Id)
	}
	return url, nil
}
