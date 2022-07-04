package rspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileService struct {
	BaseService
}

func (fs *FileService) filesUrl() string {
	return fs.BaseUrl.String() + "/files"
}

// Paginated listing of Files. Optionally the listing can be filtered by a media type
// of 'document', image', or 'av'
func (fs *FileService) Files(config RecordListingConfig, mediaType string) (*FileList, error) {

	var validMediaTypes = []string{"document", "av", "image"}
	params := config.toParams()
	if len(mediaType) > 0 {
		if ok := validateArrayContains(validMediaTypes, []string{mediaType}); !ok {
			return nil, errors.New("Invalid media type: Must be one of " + strings.Join(validMediaTypes, ","))
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
// Returns either a FileInfo of the created file or an error if operation did not succeed.
func (fs *FileService) UploadFile(config FileUploadConfig) (*FileInfo, error) {
	return fs._doUpload(config.FilePath, 0, config.Caption, config.FolderId)
}

// UploadFileNewVersion replaces the RSpace file of the given ID with the new file.
// The new version can have a different name but must be same filetype (i.e. have the same suffix)
func (fs *FileService) UploadFileNewVersion(path string, fileToReplaceId int) (*FileInfo, error) {
	return fs._doUpload(path, fileToReplaceId, "", 0)
}

func (fs *FileService) _doUpload(path string, fileToReplaceId int, caption string, folderId int) (*FileInfo, error) {
	if fileToReplaceId < 0 {
		return nil, fmt.Errorf("fileToReplaceId should be 0 or a real ID, not %d", fileToReplaceId)
	}
	urlS := fs.filesUrl()
	if fileToReplaceId != 0 {
		urlS = fmt.Sprintf("%s/%d/file", urlS, fileToReplaceId)
	}
	params := url.Values{}
	if len(caption) > 0 {
		params.Add("caption", caption)
	}
	if folderId > 0 {
		params.Add("folderId", strconv.Itoa(folderId))
	}
	encoded := params.Encode()
	if len(encoded) > 0 {
		urlS = urlS + "?" + encoded
	}
	resp, err := fs.doMultipart(path, urlS)
	if err != nil {
		return nil, err
	}
	result := &FileInfo{}
	Unmarshal(resp, result)
	return result, nil
}

func (bs *BaseService) doMultipart(path string, url string) (*http.Response, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))

	if err != nil {
		return nil, err
	}
	_, _ = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	hc := HttpClientNew(10)
	retry := NewResilientClient(hc)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	bs.addAuthHeader(req)
	resp, err := retry.Do(req)
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	return resp, nil
}

// DownloadFile retrieves the given file from RSpace and downloads to the specified directory  on local machine, which must be a writable file.
// Returns the FileInfo metadata for the downloaded file
func (fs *FileService) DownloadFile(fileId int, outDir string) (*FileInfo, error) {
	downloadUrl := fmt.Sprintf("%s/%d/file", fs.filesUrl(), fileId)
	info, err := fs.FileById(fileId)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(outDir, info.GetName())
	out, _ := os.Create(path)
	err = fs.doGetToFile(downloadUrl, out)
	if err != nil {
		return nil, err
	}
	return info, nil
}
