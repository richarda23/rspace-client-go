package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func GetStatus() *Status {
	statusStr := DoGet(BASE_URL + "/status")
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	fmt.Println(res)
	return &res
}

// Paginated listing of Documents
func Documents(config RecordListingConfig) *DocumentList {
	url := DOCUMENTS_URL + "?pageSize=" + strconv.Itoa(config.PageSize) + "&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson := DoGet(url)
	var result = DocumentList{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

//DocumentById retrieves full document content
func DocumentById(docId int) *Document {
	url := fmt.Sprintf("%s/%d", DOCUMENTS_URL, docId)
	docJson := DoGet(url)
	fmt.Println(docJson)
	var result = Document{}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

func DocumentNew(post *DocumentPost) *DocumentInfo {
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", DOCUMENTS_URL, bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}
func NewBasicDocumentWithContent(name string, tags string, contentHtml string) *DocumentInfo {
	post := BasicPost(name, tags)
	content := FieldContent{contentHtml}
	fields := make([]FieldContent, 1)
	fields[0] = content
	post.Fields = fields
	return doPostCreateDocument(post)
}

func NewEmptyBasicDocument(name string, tags string) *DocumentInfo {
	post := BasicPost(name, tags)
	return doPostCreateDocument(post)
}

func doPostCreateDocument(postData *DocumentPost) *DocumentInfo {
	hc := http.Client{}
	formData, _ := json.Marshal(postData)
	req, err := http.NewRequest("POST", DOCUMENTS_URL, bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	result := &DocumentInfo{}
	Unmarshal(resp, result)
	return result
}

func main() {
	fmt.Println("Starting the application...")
	response, err := http.Get("https://httpbin.org/ip")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	jsonData := map[string]string{"firstname": "Nic", "lastname": "Raboy"}
	jsonValue, _ := json.Marshal(jsonData)
	response, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
	fmt.Println("Terminating the application...")
	GetStatus()
	//    Documents()
}
