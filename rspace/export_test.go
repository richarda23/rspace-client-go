package rspace

import (
	"fmt"
	"testing"
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
	job, err := webClient.Export(post, true)
	if err != nil {
		Log.Info(err)
		t.Fatalf("Error creating job %s", err)
	}
	assertStringEquals(t, "COMPLETED", job.Status, "")
	url := job.DownloadLink()
	assertNotNil(t, url, "download link should be present")

	// test submit, non-blocking
	job, err = webClient.Export(post, false)
	assertStringEquals(t, "STARTING", job.Status, "")
}

func TestMakeExportUrl(t *testing.T) {
	es := ExportService{}
	post := ExportPost{XML_FORMAT, USER_EXPORT_SCOPE, 5}
	url := es.makeUrl(post, "/export")
	assertStringEquals(t, "/export/xml/user/5", url, "")

	post = NewExportPost()
	url = es.makeUrl(post, "/export")
	assertStringEquals(t, "/export/html/user", url, "")
}
