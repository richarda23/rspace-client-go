package rspace

import (
	"encoding/json"
	"time"
	"net/url"
	"strconv"
	"io/ioutil"
//	"io/ioutil"
)

type ImportService struct {
	BaseService
}

func (fs *ImportService) importUrl() string {
	return fs.BaseUrl.String() + "/import"
}

// ImportWord takes an MSWord or rich text file and imports it as a native RSpace document.
// If `folderId` is not specified, the document will be created in RSpace Home folder
// If `imageFolderId` is not specified, any images embedded in the original document will be put in the `ApiInbox`
// folder of the Image Gallery.
func (fs *ImportService) ImportWord(path string, folderId int, imageFolderId int) (*DocumentInfo, error) {
	time.Sleep(fs.Delay)
	urlStr := fs.importUrl() + "/word"
	params := url.Values{}
	if folderId > 0 {
		params.Add("folderId", strconv.Itoa(folderId))
	}
	if imageFolderId > 0 {
		params.Add("imageFolderId", strconv.Itoa(imageFolderId))
	}
	paramStr := params.Encode()
	if len(paramStr)> 0 {
		urlStr = urlStr + "?" + paramStr
	}
	resp, err := fs.doMultipart(path, urlStr)
	if err != nil {
	 return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err2 := testResponseForError(data, resp); err2 != nil {
		return nil, err2
	}
	result := &DocumentInfo{}
	json.Unmarshal(data, result)
	return result, nil
}