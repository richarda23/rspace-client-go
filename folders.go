package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func foldersUrl () string {
	return getenv (BASE_URL_ENV_NAME) + "/folders"
}

// FolderTree produces paginated listing of items in folder. If folderId is 0 then Home Folder is lister
func FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) *FolderList {
	url := foldersUrl() + "/tree"
	if folderId != 0 {
		url = url +"/" + strconv.Itoa(folderId)
	}
	url = url  + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	if len(typesToInclude) > 0 {
		url = url + "&typesToInclude=" + strings.Join(typesToInclude, ",")
	}
	docJson := DoGet(url)
	var result = FolderList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

//DocumentById retrieves full information about the folder
func FolderById(folderId int) *Folder {
	url := fmt.Sprintf("%s/%d", foldersUrl(), folderId)
	docJson := DoGet(url)
	var result = Folder{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}
// FolderNew creates a new folder or notebook with  the given name.
// If a parentFolderId is specified then the folder is created in that folder
func FolderNew(post *FolderPost) (*Folder, error) {
	var formData []byte
	if post.ParentFolderId == 0 {
		noIdPost := struct {
			Name     string `json:"name"`
			Notebook bool   `json:"notebook"`
		}{
			post.Name,
			post.IsNotebook,
		}
		fmt.Println(noIdPost)
		formData, _ = json.Marshal(&noIdPost)
	} else {
		formData, _ = json.Marshal(post)
	}
	hc := http.Client{}
	req, err := http.NewRequest("POST", foldersUrl(), bytes.NewBuffer(formData))
	if err != nil {
		return nil, err
	}
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		Log.Error(err)
	}
	result := &Folder{}
	Unmarshal(resp, result)

	fmt.Println(result)
	return result, nil
}
