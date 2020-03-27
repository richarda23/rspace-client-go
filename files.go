package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func filesUrl() string {
	return getenv (BASE_URL_ENV_NAME) + "/files"
}

// Paginated listing of Files
func Files(config RecordListingConfig) *FileList {
	url := filesUrl() + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson := DoGet(url)
	var result = FileList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

// FileById retrieves file information for a single File
func FileById(fileId int) *FileInfo {
	url := fmt.Sprintf("%s/%d", filesUrl(), fileId)
	docJson := DoGet(url)
	var result = FileInfo{}
	json.Unmarshal([]byte(docJson), &result)
	fmt.Println(docJson)
	return &result
}

// UploadFile uploads the file specified to the 'ApiInbox' subfolder of the
// appropriate Gallery section
// Panics if file cannot be read.
// Returns either a FileInfo of the created file or an error if operation did not succeed.
func UploadFile(path string) (*FileInfo, error) {
	return _doUpload(path, 0)
}

// UploadFileNewVersion replaces the RSpace file of the given ID with the new file.
// The new version can have a different name but must be same filetype (i.e. have the same suffix)
func UploadFileNewVersion(path string, fileToReplaceId int) (*FileInfo, error) {
	return _doUpload(path, fileToReplaceId)
}

func _doUpload(path string, fileToReplaceId int) (*FileInfo, error) {
	if fileToReplaceId < 0 {
		return nil, fmt.Errorf("fileToReplaceId should be 0 or a real ID, not %d", fileToReplaceId)
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	hc := http.Client{}
	url := filesUrl()
	if fileToReplaceId != 0 {
		url = fmt.Sprintf("%s/%d/file", url, fileToReplaceId)
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	AddAuthHeader(req)
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result := &FileInfo{}
	Unmarshal(resp, result)
	return result, nil
}

// DownloadFile retrieves the given file from RSpace and downloads to the specified file location on local machine, which must be a writable file.
func DownloadFile(fileId int, outFile string) {
	url := fmt.Sprintf("%s/%d/file", filesUrl(), fileId)
	err := DoGetToFile(url, outFile)
	if err != nil {
		log.Fatalln(err)
	}
}
