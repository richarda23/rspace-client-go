package rspace

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	APIKEY   = "4FmuGC6OCVlW8QqNNz448PEMCutJtgBL"
	BASE_URL = "https://demos.researchspace.com/api/v1"
	DOCUMENTS_URL=BASE_URL+"/documents"
	FILES_URL=BASE_URL+"/files"
)

func BasicPost (name string, tags string) *DocumentPost {
	post := DocumentPost{}
	post.Name = name
	if len(tags) > 0 {
	  post.Tags=tags
	}
	return &post
}

func Unmarshal(resp *http.Response, result interface{} )  {
	data, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, result)
}
func AddAuthHeader (req *http.Request) {
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

func Abbreviate(toAbbreviate string, maxLen int) string {
	if len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-4)] + "..."
	}
	return toAbbreviate
}

func DoGet(url string) string {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	AddAuthHeader(req)
	resp, _ := client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	respStr := string(data)
	return respStr
}

