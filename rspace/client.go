package rspace

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func getenv(envname string) string {
	return os.Getenv(envname)
}

func validateArrayContains(validTerms []string, toTest []string) bool {
	for _, term := range toTest {
		seen := false
		for _, v := range validTerms {
			if v == term {
				seen = true
			}
		}
		if !seen {
			return false
		}
	}
	return true
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
	//fmt.Println(string(data))
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


// Abbreviate truncates a string to maximum length `maxLen`, including
// 3 ellipsis characters.
func abbreviate(toAbbreviate string, maxLen int) string {
	if len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-4)] + "..."
	}
	return toAbbreviate
}

func testResponseForError(data []byte, resp *http.Response) *RSpaceError {
	if resp.StatusCode >= 400 {
		rspaceError := &RSpaceError{}
		json.Unmarshal(data, rspaceError)
		return rspaceError
	}
	return nil
}