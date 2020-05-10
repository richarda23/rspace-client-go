package rspace

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type FolderService struct {
	BaseService
}

func (fs *FolderService) foldersUrl() string {
	return fs.BaseUrl.String() + "/folders"
}

// FolderTree produces paginated listing of items in folder. If folderId is 0 then Home Folder is lister
func (fs *FolderService) FolderTree(config RecordListingConfig, folderId int, typesToInclude []string) (*FolderList, error) {
	time.Sleep(fs.Delay)
	urlStr := fs.foldersUrl() + "/tree"
	if folderId != 0 {
		urlStr = urlStr + "/" + strconv.Itoa(folderId)
	}
	params := config.toParams()
	if len(typesToInclude) > 0 {
		params.Add("typesToInclude", strings.Join(typesToInclude, ","))
	}
	encoded := params.Encode()
	if len(encoded) > 0 {
		urlStr = urlStr + "?" + encoded
	}
	//fmt.Println(url)
	data, err := fs.doGet(urlStr)
	if err != nil {
		return nil, err
	}
	var result = FolderList{}
	json.Unmarshal(data, &result)
	return &result, nil
}

//DocumentById retrieves full information about the folder
func (fs *FolderService) FolderById(folderId int) (*Folder, error) {
	time.Sleep(fs.Delay)
	url := fmt.Sprintf("%s/%d", fs.foldersUrl(), folderId)
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = Folder{}
	json.Unmarshal(data, &result)
	return &result, err
}

// DeleteFolder attempts to delete the folder or noteboon with the specified ID
func (fs *FolderService) DeleteFolder(folderId int) (bool, error) {
	time.Sleep(fs.Delay)
	url := fmt.Sprintf("%s/%d", fs.foldersUrl(), folderId)
	resp, err := fs.doDelete(url)
	if resp == false {
		return false, err
	} else {
		return true, nil
	}
}

// FolderNew creates a new folder or notebook with  the given name.
// If a parentFolderId is specified then the folder is created in that folder
func (fs *FolderService) FolderNew(post *FolderPost) (*Folder, error) {
	time.Sleep(fs.Delay)
	data, err := fs.doPostJsonBody(post, fs.foldersUrl())

	if err != nil {
		return nil, err
	}
	result := &Folder{}
	json.Unmarshal(data, result)
	return result, nil
}
