package rspace

import (
	"fmt"
	"testing"
	"time"
)
var folderService *FolderService = &FolderService{
	BaseService:BaseService{
		Delay:time.Duration(100) * time.Millisecond}}

func TestNewFolderGetFolder(t *testing.T) {
	post := FolderPost{}
	post.Name = "f1"
	post.IsNotebook = true
	got, err := folderService.FolderNew(&post)
	Log.Info(Marshal(got))
	if err != nil {
		Log.Error(err)
	}

	if got.Name != "f1" {
		fail(t, fmt.Sprintf("expected name %s  but was %s", "f1", got.Name))
	}
	folder,e := folderService.FolderById(got.Id)
	if e != nil {
		Log.Error(e)
	}
	Log.Info(Marshal(got))
	if folder.IsNotebook == true {
		fail(t, fmt.Sprintf("expected folder, not notebook"))
	}
	if folder.Id != got.Id {
		fail(t, fmt.Sprintf("expected ID = %d, but was %d", got.Id, folder.Id))
	}
}
func TestListFolderTree(t *testing.T) {
	cfg := NewRecordListingConfig()
	types := make([]string, 1)
	types[0]="notebook"
	// to do fix 'types' usage
	result,e := folderService.FolderTree(cfg, 0, types)
	if e != nil {
		Log.Error(e)
	}
	for _, v := range result.Records {
		if v.Type != "NOTEBOOK" {
			fail (t, fmt.Sprintf("Folder listing should be notebooks only"))
		}
	}
}
func TestErrorHandling (t *testing.T) {
	folder,e := folderService.FolderById(-233)
	if folder != nil {
			fail (t, fmt.Sprintf("Should have invoked an error"))
	}
	if e == nil {
			fail (t, fmt.Sprintf("Error object should not be nil"))
	}
	Log.Info(e)
}
