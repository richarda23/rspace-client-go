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
	"time"
	"strings"
	"errors"
)

type FileService struct {
	BaseService
}

func (fs *FileService) filesUrl() string {
	return fs.BaseUrl.String()+ "/files"
}


// Paginated listing of Files. Optionally the listing can be filtered by a media type
// of 'document', image', or 'av'
func (fs *FileService) Files(config RecordListingConfig, mediaType string) (*FileList, error) {
	time.Sleep(fs.Delay)
	var validMediaTypes  = []string{"document", "av", "image"}
	params := config.toParams()
	if len(mediaType) > 0 {
		if ok := validateArrayContains(validMediaTypes, []string{mediaType}); !ok {
			return nil, errors.New("Invalid media type: Must be one of " + strings.Join(validMediaTypes,","))
		}
		params.Add("mediaType", mediaType)
	}
	
	encoded := params.Encode()
	url := fs.filesUrl() + "?" + encoded
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = FileList{}
	json.Unmarshal(data, &result)
	return &result, nil
}

// FileById retrieves file information for a single File
func (fs *FileService) FileById(fileId int) (*FileInfo, error) {
	time.Sleep(fs.Delay)
	url := fmt.Sprintf("%s/%d", fs.filesUrl(), fileId)
	data, err := fs.doGet(url)
	if err != nil {
		return nil, err
	}
	var result = FileInfo{}
	json.Unmarshal(data, &result)
	return &result, nil
}

// UploadFile uploads the file specified to the 'ApiInbox' subfolder of the
// appropriate Gallery section
// Panics if file cannot be read.
// Returns either a FileInfo of the created file or an error if operation did not succeed.
func (fs *FileService) UploadFile(path string) (*FileInfo, error) {
	time.Sleep(fs.Delay)
	return fs._doUpload(path, 0)
}

// UploadFileNewVersion replaces the RSpace file of the given ID with the new file.
// The new version can have a different name but must be same filetype (i.e. have the same suffix)
func (fs *FileService) UploadFileNewVersion(path string, fileToReplaceId int) (*FileInfo, error) {
	time.Sleep(fs.Delay)
	return fs._doUpload(path, fileToReplaceId)
}

func (fs *FileService) _doUpload(path string, fileToReplaceId int) (*FileInfo, error) {
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
	url := fs.filesUrl()
	if fileToReplaceId != 0 {
		url = fmt.Sprintf("%s/%d/file", url, fileToReplaceId)
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	fs.addAuthHeader(req)
	resp, err := hc.Do(req)
	if err != nil {
		Log.Error(err)
	}
	result := &FileInfo{}
	Unmarshal(resp, result)
	return result, nil
}

// DownloadFile retrieves the given file from RSpace and downloads to the specified directory  on local machine, which must be a writable file.
// Returns the FileInfo metadata for the downloaded file
func (fs *FileService) DownloadFile(fileId int, outDir string) (*FileInfo, error) {
	downloadUrl := fmt.Sprintf("%s/%d/file", fs.filesUrl(), fileId)
	info,err:=fs.FileById(fileId)
	path := filepath.Join(outDir, info.GetName())
	if err != nil {
		return nil, err
	}
	err = fs.doGetToFile(downloadUrl, path)
	if err != nil {
		return nil, err
	}
	return info, nil
}
