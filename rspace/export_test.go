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
