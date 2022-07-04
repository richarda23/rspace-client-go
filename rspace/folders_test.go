package rspace

import (
	"fmt"
	"testing"
)

func TestNewUntitledFolder(t *testing.T) {
	post := FolderPost{}
	got, _ := webClient.FolderNew(&post)
	assertNotNil(t, got.Name, "")
	assertFalse(t, got.IsNotebook, "should be folder, not notebook")
}
func TestNewFolderGetFolder(t *testing.T) {
	post := FolderPost{}
	post.Name = "f1"
	post.IsNotebook = true
	got, err := webClient.FolderNew(&post)
	if err != nil {
		Log.Error(err)
	}

	if got.Name != "f1" {
		fail(t, fmt.Sprintf("expected name %s  but was %s", "f1", got.Name))
	}
	folder, _ := webClient.FolderById(got.Id)
	fmt.Println(folder)
	if !folder.IsNotebook {
		fail(t, "expected folder, not notebook")
	}
	if folder.Id != got.Id {
		fail(t, fmt.Sprintf("expected ID = %d, but was %d", got.Id, folder.Id))
	}
	// delete
	deletionResult, _ := webClient.DeleteFolder(folder.Id)
	assertTrue(t, deletionResult, "deletion of folder did not succeed")
	// now get by ID  should fail
	_, e2 := webClient.FolderById(got.Id)
	rsErr, ok := e2.(*RSpaceError)
	assertTrue(t, ok, "could not convert to RSpace error")
	assertIntEquals(t, 401, rsErr.HttpCode, "")
}
func TestListFolderTree(t *testing.T) {
	cfg := NewRecordListingConfig()
	types := make([]string, 1)
	types[0] = "notebook"
	// to do fix 'types' usage
	result, e := webClient.FolderTree(cfg, 0, types)
	if e != nil {
		Log.Error(e)
	}
	for _, v := range result.Records {
		if v.Type != "NOTEBOOK" {
			fail(t, "Folder listing should be notebooks only")
		}
	}
}
func TestErrorHandling(t *testing.T) {
	folder, e := webClient.FolderById(-233)
	if folder != nil {
		fail(t, "Should have invoked an error")
	}
	if e == nil {
		fail(t, "Error object should not be nil")
	}
	Log.Info(e)
}
