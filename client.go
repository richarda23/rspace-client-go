package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

const (
	APIKEY   = "4FmuGC6OCVlW8QqNNz448PEMCutJtgBL"
	BASE_URL = "https://demos.researchspace.com/api/v1"
)
type Status struct {
	Message       string `json: "message"`
	RSpaceVersion string `json: "rspaceVersion"`
}

type RecordListingConfig struct {
    SortOrder string
    PageSize int
    PageNumber int
    OrderBy string
    Quiet bool
}
func New () RecordListingConfig {
	return RecordListingConfig{
         PageSize:20,
	 OrderBy:"lastModified",
	 PageNumber:1,
	 SortOrder:"desc",
	 Quiet: false,
	}
}


func GetStatus() *Status {
	statusStr := doGet(BASE_URL + "/status")
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	fmt.Println(res)
	return &res
}

func Documents(config RecordListingConfig) {
	fmt.Println( config)
        url := BASE_URL + "/documents?pageSize=" + strconv.Itoa(config.PageSize) +"&pageNumber=" + strconv.Itoa(config.PageNumber)
	fmt.Println("url is " + url)
	docJson := doGet(url)
	var result map[string]interface{}
	json.Unmarshal([]byte(docJson), &result)
	docs := result["documents"].([]interface{})

	for _, value := range docs {
		item := value.(map[string]interface{})
		// Each value is an interface{} type, that is type asserted as a string
		id := int(item["id"].(float64))
		name := abbreviate(item["name"].(string), 30)
		t, _ := time.Parse(time.RFC3339Nano, item["lastModified"].(string))
		lm := t.Format(time.RFC3339)
		if config.Quiet {
			fmt.Printf("%-10d\n", id)
		} else {
			fmt.Printf("%-10d%-30s%-20s\n", id, name, lm)
		}
	}

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
	req.Header.Add("apiKey", APIKEY)
	resp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	respStr := string(data)
	fmt.Println(string(respStr))
	return string(respStr)
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
