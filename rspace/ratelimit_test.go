package rspace

import (
	"net/http"
	"testing"
)

func TestRateLimitParse(t *testing.T) {
	resp := createValidRateLimitHeaders()
	var rld RateLimitData = NewRateLimitData(&resp)
	assertIntEquals(t, 100, rld.RateLimit, "")
	assertIntEquals(t, 50, rld.Remaining, "")
	assertIntEquals(t, 5, rld.MinWaitIntervalMillis, "")

	// unknown value set to -100
	resp.Header.Del(RATE_LIMIT_HDR)
	resp.Header.Add(RATE_LIMIT_HDR, "")
	rld = NewRateLimitData(&resp)
	assertIntEquals(t, -100, rld.RateLimit, "")
}

//todo create ttest header
func createValidRateLimitHeaders() http.Response {
	resp := http.Response{}
	hdr := make(map[string][]string)
	h := http.Header(hdr)
	h.Add(RATE_LIMIT_HDR, "100")
	h.Add(RATE_LIMIT_REMAINING_HDR, "50")
	h.Add(RATE_LIMIT_MIN_WAIT_HDR, "5")
	resp.Header = h
	return resp
}
