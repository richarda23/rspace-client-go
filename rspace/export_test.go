package rspace

import (
	"fmt"
	"testing"
	"time"
)

func TestExportScope(t *testing.T) {
	scoper := USER_EXPORT_SCOPE
	assertStringEquals(t, "user", scoper.String(), "")
}

func TestExportFormat(t *testing.T) {
	format := XML_FORMAT
	assertStringEquals(t, "xml", format.String(), "")
}

func TestJobDownloadLink(t *testing.T) {
	job := Job{}
	// download link is initially nil
	fmt.Println(job.DownloadLink())
	if job.DownloadLink() != nil {
		t.Fatalf("download link wasn't nil")
	}
}

func TestDoExport(t *testing.T) {
	// export a new user
	userPost := createRandomUser(Pi)
	newUser, _ := webClient.UserNew(userPost)
	post := NewExportPost()
	post.Id = newUser.Id
	job, err := webClient.Export(post, true, func(string) {})
	if err != nil {
		Log.Info(err)
		t.Fatalf("Error creating job %s", err)
	}
	assertStringEquals(t, "COMPLETED", job.Status, "")
	url := job.DownloadLink()
	assertNotNil(t, url, "download link should be present")

	// now submit selection
	doc0, _ := webClient.Documents(NewRecordListingConfig())
	id0 := doc0.Documents[0].Id
	post.Scope = SELECTION_EXPORT_SCOPE
	post.Id = 0
	post.ItemIds = []int{id0}
	job, err = webClient.Export(post, true, func(string) {})
	assertStringEquals(t, "COMPLETED", job.Status, "")

	// test submit, non-blocking
	job, err = webClient.Export(post, false, func(string) {})
	assertStringEquals(t, "STARTING", job.Status, "")
}

var messages = make([]string, 1)

func reporter(s string) {
	messages = append(messages, s)
}
func TestCalculateRemainingTime(t *testing.T) {
	start := time.Now().Add(time.Second * -100)
	//50% done in 100s means 50% left to do - another 100s
	// poll with 1/5th interval, expected interval = 100/5= 20
	dur, _ := calculateSleepTime(50, start, reporter)
	for _, item := range messages {
		println(item)
	}
	assertDurationEquals(t, time.Duration(time.Second*20), *dur, "")

	// max of 60 seconds re-polling time
	dur, _ = calculateSleepTime(0.01, start, reporter)
	assertDurationEquals(t, time.Duration(time.Second*60), *dur, "")

	// minimum of 3 second interval
	dur, _ = calculateSleepTime(99.999, start, reporter)
	assertDurationEquals(t, time.Duration(time.Second*3), *dur, "")

	//error if progress is 0 - can't calculate
	_, err := calculateSleepTime(0.0, start, reporter)
	assertNotNil(t, err, "")

}

func TestMakeExportUrl(t *testing.T) {
	es := ExportService{}
	post := ExportPost{XML_FORMAT, USER_EXPORT_SCOPE, 5, []int{}, 1}
	url, _ := es.makeUrl(post, "/export")
	assertStringEquals(t, "/export/xml/user/5", url, "")

	post = NewExportPost()
	url, _ = es.makeUrl(post, "/export")
	assertStringEquals(t, "/export/html/user", url, "")

	post = NewExportPost()
	post.ItemIds = []int{1, 2, 3}
	post.Scope = SELECTION_EXPORT_SCOPE

	url, _ = es.makeUrl(post, "/export")

	assertStringEquals(t, "/export/html/selection?selections=1,2,3&maxLinkLevel=1", url, "")

}
