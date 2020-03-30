package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileService struct {
	BaseService
}

func filesUrl() string {
	return getenv(BASE_URL_ENV_NAME) + "/files"
}

// Paginated listing of Files
func (fs *FileService) Files(config RecordListingConfig) (*FileList, error) {
	time.Sleep(fs.Delay)
	url := filesUrl() + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = FileList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result, nil
}

// FileById retrieves file information for a single File
func (fs *FileService) FileById(fileId int) (*FileInfo, error) {
	time.Sleep(fs.Delay)
	url := fmt.Sprintf("%s/%d", filesUrl(), fileId)
	docJson, err := DoGet(url)
	if err != nil {
		return nil, err
	}
	var result = FileInfo{}
	json.Unmarshal([]byte(docJson), &result)
	fmt.Println(docJson)
	return &result, nil
}

// UploadFile uploads the file specified to the 'ApiInbox' subfolder of the
// appropriate Gallery section
// Panics if file cannot be read.
// Returns either a FileInfo of the created file or an error if operation did not succeed.
func (fs *FileService) UploadFile(path string) (*FileInfo, error) {
	time.Sleep(fs.Delay)
	return _doUpload(path, 0)
}

// UploadFileNewVersion replaces the RSpace file of the given ID with the new file.
// The new version can have a different name but must be same filetype (i.e. have the same suffix)
func (fs *FileService) UploadFileNewVersion(path string, fileToReplaceId int) (*FileInfo, error) {
	time.Sleep(fs.Delay)
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
		Log.Error(err)
	}
	result := &FileInfo{}
	Unmarshal(resp, result)
	return result, nil
}

// DownloadFile retrieves the given file from RSpace and downloads to the specified file location on local machine, which must be a writable file.
func (fs *FileService) DownloadFile(fileId int, outFile string) {
	url := fmt.Sprintf("%s/%d/file", filesUrl(), fileId)
	err := DoGetToFile(url, outFile)
	if err != nil {
		Log.Error(err)
	}
}
