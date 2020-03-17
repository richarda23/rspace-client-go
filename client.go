package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"log"
)

const (
	APIKEY   = "4FmuGC6OCVlW8QqNNz448PEMCutJtgBL"
	BASE_URL = "https://demos.researchspace.com/api/v1"
	DOCUMENTS_URL=BASE_URL+"/documents"
)


func GetStatus() *Status {
	statusStr := doGet(BASE_URL + "/status")
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	fmt.Println(res)
	return &res
}
// Paginated listing of Documents
func Documents(config RecordListingConfig) *DocumentList {
        url := DOCUMENTS_URL + "?pageSize=" + strconv.Itoa(config.PageSize) +"&pageNumber=" + strconv.Itoa(config.PageNumber)
	docJson := doGet(url)
	var result = DocumentList {}
	json.Unmarshal([]byte(docJson), &result)
	return &result
}

func DocumentNew (post *DocumentPost) *DocumentInfo {
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", DOCUMENTS_URL, bytes.NewBuffer(formData))
	addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return unmarshal(resp)
}

func NewEmptyBasicDocument (name string, tags string) *DocumentInfo {
	hc := http.Client{}
	post := DocumentPost{}
	post.Name = name
	if len(tags) > 0 {
	  post.Tags=tags
	}
	fmt.Println(post)
	formData, _ := json.Marshal(post)
	fmt.Println(string(formData))
	req, err := http.NewRequest("POST", DOCUMENTS_URL, bytes.NewBuffer(formData))
	addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return unmarshal(resp)
}
func unmarshal(resp *http.Response) *DocumentInfo {
	data, _ := ioutil.ReadAll(resp.Body)
	var result = &DocumentInfo {}
	json.Unmarshal(data, result)
	return result
}
func addAuthHeader (req *http.Request) {
	req.Header.Add("apiKey", APIKEY)
}

func PrintDocs () {

//	docs := result["documents"].([]interface{})
//	for _, value := range docs {
//		item := value.(map[string]interface{})

		// Each value is an interface{} type, that is type asserted as a string
//		id := int(item["id"].(float64))
//		name := abbreviate(item["name"].(string), 30)
//		t, _ := time.Parse(time.RFC3339Nano, item["lastModified"].(string))
//		lm := t.Format(time.RFC3339)
//		if config.Quiet {
//			fmt.Printf("%-10d\n", id)
//		} else {
//			fmt.Printf("%-10d%-30s%-20s\n", id, name, lm)
//		}
//	}

}

func abbreviate(toAbbreviate string, maxLen int) string {
	if len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-4)] + "..."
	}
	return toAbbreviate
}

func doGet(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	addAuthHeader(req)
	resp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	respStr := string(data)
	return respStr
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
