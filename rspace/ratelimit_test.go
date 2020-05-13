package rspace

import (
	"net/http"
	"testing"
)

func TestRateLimitParse(t *testing.T) {
	resp := createValidRateLimitHeaders()
	var rld RateLimitData = NewRateLimitData(&resp)
	assertIntEquals(t, 100, rld.WaitTimeMillis, "")

	// unknown value set to -100
	resp.Header.Del(RATE_LIMIT_WAIT_TIME)
	resp.Header.Add(RATE_LIMIT_WAIT_TIME, "")
	rld = NewRateLimitData(&resp)
	assertIntEquals(t, -100, rld.WaitTimeMillis, "")
}

//todo create ttest header
func createValidRateLimitHeaders() http.Response {
	resp := http.Response{}
	hdr := make(map[string][]string)
	h := http.Header(hdr)
	h.Add(RATE_LIMIT_WAIT_TIME, "100")
	resp.Header = h
	return resp
}
