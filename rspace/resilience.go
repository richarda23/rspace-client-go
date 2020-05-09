package rspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// ClientEx abstracts out the http.Client Do() method into an interface
// for decorating
type ClientEx interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryClientEx performs a fixed number of attempts to call clientEx.Do()
type RetryClientEx struct {
	retries int
	cli     ClientEx
}

// RetryClientExNew is constructor for RetryClientEx. It validates the number of retries >= 1
func RetryClientExNew(retries int, wrappedRequest ClientEx) (*RetryClientEx, error) {
	if retries < 1 {
		return nil, errors.New("Number of retries must be >= 1")
	} else {
		return &RetryClientEx{retries, wrappedRequest}, nil
	}
}

func (ex RetryClientEx) Do(req *http.Request) (*http.Response, error) {
	var currErr error
	var resp *http.Response
	for i := 0; i < ex.retries; i++ {
		Log.Infof("Retry attempt %d", i)
		resp, currErr = ex.cli.Do(req)
		// e.g. server not available
		if currErr != nil {
			return nil, currErr
		}
		if x := testResponseForError(resp); x != nil {
			//is this error worth retrying? Don't retry client error
			// unless is 429
			if x.HttpCode == 429 || x.HttpCode > 500 {
				// we have an error
				Log.Infof(" Got an error - %s , retrying", x)
			} else {
				return nil, x
			}
		} else {
			return resp, nil
		}
	}
	return nil, currErr
}

type DelayClientEx struct {
	cli ClientEx
}

func (this *DelayClientEx) Do(req *http.Request) (*http.Response, error) {

	resp, e := this.cli.Do(req)
	if e != nil {
		Log.Error(e)
		return nil, e
	}
	var rld RateLimitData = NewRateLimitData(resp)
	Log.Info(rld.String())

	if err := testResponseForError(resp); err != nil {
		if err.HttpCode == 429 {
			Log.Warningf("429 error, waiting for %d ms till next call", err.MillisTillNextCall)
			time.Sleep(time.Duration(err.MillisTillNextCall) * time.Millisecond)
		}
		return nil, err
	}
	return resp, nil

}

// testResponseForError reads response body and if error code > 400
// will construct an RSpaceError.
// It will reset the response so that it can be read again
func testResponseForError(resp *http.Response) *RSpaceError {
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	if resp.StatusCode >= 400 {
		rspaceError := &RSpaceError{}
		json.Unmarshal(data, rspaceError)
		if resp.StatusCode == 429 {
			// set min wait time from header
			// TODO change
			rspaceError.MillisTillNextCall = 1000
		}
		return rspaceError
	}
	return nil
}

//RSpaceError encapsulates server or client side errors leading to a request being rejected.
type RSpaceError struct {
	Status             string
	HttpCode           int
	InternalCode       int
	Message            string
	Errors             []string
	Timestamp          string `json:"iso8601Timestamp"`
	MillisTillNextCall int
}

func (f *RSpaceError) CreatedTime() (time.Time, error) {
	return parseTimestamp(f.Timestamp)
}

func (rsError *RSpaceError) String() string {
	if rsError.HttpCode >= 400 && rsError.HttpCode < 500 {
		return formatErrorMsg(rsError, "Client")
	} else if rsError.HttpCode > 500 {
		return formatErrorMsg(rsError, "Server")
	} else {
		return formatErrorMsg(rsError, "Unknown")
	}
}

func (rsError *RSpaceError) Error() string {
	return rsError.String()
}

func formatErrorMsg(rsError *RSpaceError, errType string) string {
	concatenateErrM := strings.Join(rsError.Errors, "\n")
	rc := fmt.Sprintf("%s error:httpCode=%d, status=%s, internalCode=%d, timestamp=%s,  message=%s\nErrors: %s",
		errType, rsError.HttpCode, rsError.Status, rsError.InternalCode, rsError.Timestamp, rsError.Message, concatenateErrM)
	return rc
}
