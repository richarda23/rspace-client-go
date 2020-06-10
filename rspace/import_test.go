package rspace

import (
	"fmt"
	"testing"
)

const (
	TESTFILE_IMPORT = "testdata/BrDUlabelling.doc"
)

func TestImportWordFile(t *testing.T) {
	//upload into any folder

	got, err := webClient.ImportWord(TESTFILE_IMPORT, 0, 0)
	if err != nil {
		fmt.Println(err)
	} else {
		assertStringEquals(t, "BrDUlabelling", got.Name, "")
	}
	//create a new folder to import into
	got2, _ := webClient.FolderNew(&FolderPost{})
	newFolderId := got2.Id
	got3, _ := webClient.ImportWord(TESTFILE_IMPORT, newFolderId, 0)
	assertIntEquals(t, newFolderId, got3.ParentFolderId, "")

	//check error handling
	_, err4 := webClient.ImportWord(TESTFILE_IMPORT, 1234567, 0)
	assertNotNil(t, err4, "")
	rs_err := err4.(*RSpaceError)
	assertIntEquals(t, 401, rs_err.HttpCode, "expect 401 auth erro as folder 1234567 doesn't exist")

}
