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

// var filesService *FileService = &FileService{
// 	BaseService: BaseService{
// 		Delay: time.Duration(100) * time.Millisecond}}

func TestFileList(t *testing.T) {
	cfg := NewRecordListingConfig()

	_, err1 := webClient.Files(cfg, "invalidMediaType")
	assertNotNil(t, err1, "invalid media type should cause error")
	got, err := webClient.Files(cfg, "")
	if err != nil {
		Log.Error(err)
	}
	if got.TotalHits <= 1 {
		fail(t, fmt.Sprintf("Expected hits > 1 but was %d", got.TotalHits))
	}
	fmt.Printf("Got %d hits\n", got.TotalHits)
	id := got.Files[0].Id

	file, _ := webClient.FileById(id)
	fmt.Println(file.Id)
}
func nameFromPath(path string) string {
	return strings.Split(path, "/")[1]
}
func TestFileReplace(t *testing.T) {
	cfg := FileUploadConfig{FilePath: TESTFILEUPLOAD}
	got, err := webClient.UploadFile(cfg)
	fmt.Println(got)
	fmt.Printf("uploaded id of file to replace is is %d", got.Id)
	got, err = webClient.UploadFileNewVersion(TESTFILEUPDATE, got.Id)
	if err != nil {
		fmt.Println(err)
	}
	if got.Name != nameFromPath(TESTFILEUPDATE) {
		fail(t, fmt.Sprintf("Name was %s", got.Name))
	}
}
func TestFileUpload(t *testing.T) {
	cfg := FileUploadConfig{Caption: "Caption",
		FilePath: TESTFILEUPLOAD}
	got, err := webClient.UploadFile(cfg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Uploaded file Id is %d", got.Id)
	}
	if got.Name != nameFromPath(TESTFILEUPLOAD) {
		fail(t, fmt.Sprintf("expected name %s  but was %s", nameFromPath(TESTFILEUPLOAD), got.Name))
	}
	outfile := fmt.Sprintf("/tmp/%s", got.Name)
	webClient.Download(got.Id, outfile)
}
