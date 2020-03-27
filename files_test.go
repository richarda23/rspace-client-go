package rspace

import (
	"fmt"
	"strings"
	"testing"
)

const (
	TESTFILEUPLOAD = "testdata/ServerSetupCentOS.md"
	TESTFILEUPDATE = "testdata/RSpaceConfiguration.md"
)

func TestFileList(t *testing.T) {
	cfg := NewRecordListingConfig()
	got := Files(cfg)
	if got.TotalHits <= 1 {
		fail(t, fmt.Sprintf("Expected hits > 1 but was %d", got.TotalHits))
	}
	id := got.Files[0].Id

	file := FileById(id)
	fmt.Println(file.Id)
}
func nameFromPath(path string) string {
	return strings.Split(path, "/")[1]
}
func TestFileReplace(t *testing.T) {
	got, err := UploadFile(TESTFILEUPLOAD)
	fmt.Println(got)
	fmt.Printf("uploaded id of file to replace is is %d", got.Id)
	got, err = UploadFileNewVersion(TESTFILEUPDATE, got.Id)
	if err != nil {
		fmt.Println(err)
	}
	if got.Name != nameFromPath(TESTFILEUPDATE) {
		fail(t, fmt.Sprintf("Name was %s", got.Name))
	}
}
func TestFileUpload(t *testing.T) {
	got, err := UploadFile(TESTFILEUPLOAD)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Uploaded file Id is %d", got.Id)
	}
	if got.Name != nameFromPath(TESTFILEUPLOAD) {
		fail(t, fmt.Sprintf("expected name %s  but was %s", nameFromPath(TESTFILEUPLOAD), got.Name))
	}
	outfile := fmt.Sprintf("/tmp/%s", got.Name)
	DownloadFile(got.Id, outfile)
}