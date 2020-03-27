package rspace

import (
	"fmt"
	"testing"
)

func TestNewFolderGetFolder(t *testing.T) {
	post := FolderPost{}
	post.Name = "f1"
	post.IsNotebook = true
	got, err := FolderNew(&post)
	Log.Info(Marshal(got))
	if err != nil {
		Log.Error(err)
	}

	if got.Name != "f1" {
		fail(t, fmt.Sprintf("expected name %s  but was %s", "f1", got.Name))
	}
	folder := FolderById(got.Id)
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
	result := FolderTree(cfg, 0, types)
	for _, v := range result.Records {
		if v.Type != "NOTEBOOK" {
			fail (t, fmt.Sprintf("Folder listing should be notebooks only"))
		}
	}
}
