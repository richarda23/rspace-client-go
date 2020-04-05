package rspace

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

type BaseService struct {
	Delay time.Duration
}

func getenv(envname string) string {
	return os.Getenv(envname)
}

func BasicPost(name string, tags string) *DocumentPost {
	post := DocumentPost{}
	post.Name = name
	if len(tags) > 0 {
		post.Tags = tags
	}
	return &post
}

func Unmarshal(resp *http.Response, result interface{}) {
	data, _ := ioutil.ReadAll(resp.Body)
	if data != nil {
		json.Unmarshal(data, result)
	} else {
		Log.Error("Error parsing result")
	}
}

func Marshal(anything interface{}) string {
	bytes, _ := json.Marshal(anything)
	return string(bytes)
}

func AddAuthHeader(req *http.Request) {
	req.Header.Add("apiKey", os.Getenv(APIKEY_ENV_NAME))
}

func PrintDocs() {

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

// Abbreviate truncates a string to maximum length `maxLen`, including
// 3 ellipsis characters.
func abbreviate(toAbbreviate string, maxLen int) string {
	if len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-4)] + "..."
	}
	return toAbbreviate
}

func doPostJsonBody(post interface{}, urlString string) ([]byte, error) {
	time.Sleep(ds.Delay)
	formData, _ := json.Marshal(post)
	hc := http.Client{}
	req, err := http.NewRequest("POST", urlString, bytes.NewBuffer(formData))
	AddAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err2 := testResponseForError(data, resp); err2 != nil {
		return nil, err2
	}
	return data, nil
}

//DoGet makes an authenticated API request to a URL expecting a string
// response (typically JSON)
func DoGet(url string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	AddAuthHeader(req)
	resp, e := client.Do(req)
	if e != nil {
		Log.Error(e)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := testResponseForError(data, resp); err != nil {
		return nil, err
	}
	return data, nil
}

// DoDeleteUrl attempts to delete a resource specified by the URL. If successful, returns true, else returns false, with a possible
// non-null error
func DoDelete(url string) (bool , error) {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	AddAuthHeader(req)
	resp, e := client.Do(req)
	if e != nil {
		Log.Error(e)
		return false, e
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if err := testResponseForError(data, resp); err != nil {
		return false, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	} else {
		return false, nil
	}
}
func testResponseForError(data []byte, resp *http.Response) *RSpaceError {
	if resp.StatusCode >= 400 {
		rspaceError := &RSpaceError{}
		json.Unmarshal(data, rspaceError)
		return rspaceError
	}
	return nil
}

// DoGetToFile saves the response from an HTTP GET request to the specified file.
// If the response fails or the file cannot be created returns an error.
// 'filepath' argument should be absolute path to a file. If the file exists, it will be overwritten. If it doesn't exist, it will be created.
func DoGetToFile(url string, filepath string) error {
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	AddAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	return err
}
