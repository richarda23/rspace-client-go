package rspace

import (
	"encoding/json"
	"fmt"
	"strconv"
	"os"
	"io"
	"bytes"
	"mime/multipart"
	"net/http"
	"log"
	"path/filepath"
)

// Paginated listing of Files
func Files(config RecordListingConfig) *FileList {
        url := FILES_URL + "?pageSize=" + strconv.Itoa(config.PageSize) +"&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson := DoGet(url)
	var result = FileList {}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

//FileById retrieves file information for a single File
func FileById(fileId int) *FileInfo {
        url := fmt.Sprintf("%s/%d", FILES_URL, fileId) 
	docJson := DoGet(url)
	var result = FileInfo {}
	json.Unmarshal([]byte(docJson), &result)
	fmt.Println(docJson)
	return &result
}

func UploadFile(path string) (*FileInfo, error) {
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
	req, err := http.NewRequest("POST", FILES_URL, body)
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



