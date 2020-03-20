package rspace

import (
	"testing"
	"fmt"
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
func TestFileUpload(t *testing.T) {
	got, err :=UploadFile("/home/richard/go/src/rspace/client.go")
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println(got)
	}
}


