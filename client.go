package rspace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Status struct {
	Message       string `json: "message"`
	RSpaceVersion string `json: "rspaceVersion"`
}

const (
	APIKEY   = "4FmuGC6OCVlW8QqNNz448PEMCutJtgBL"
	BASE_URL = "https://demos.researchspace.com/api/v1"
)

func GetStatus() *Status {
	statusStr := doGet(BASE_URL + "/status")
	res := Status{}
	json.Unmarshal([]byte(statusStr), &res)
	fmt.Println(res)
	return &res
}

func Documents() {
	docJson := doGet(BASE_URL + "/documents")
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
		fmt.Printf("%-10d%-30s%-20s\n", id, name, lm)
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
