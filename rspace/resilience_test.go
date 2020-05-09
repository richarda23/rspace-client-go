package rspace

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockHttpClient struct {
	numberOfInvocations int
	PassOnNthInvocation int
	// if true will produce a 500 server error response
	// if false will cause empty response
	ProduceErrorResponse bool
}

// MockClient counts the number of times it is invoked and
// will pass on the  PassOnNthInvocation-th attempt
// When it fails it produces no response at all
func (mock *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	mock.numberOfInvocations++
	if mock.numberOfInvocations >= mock.PassOnNthInvocation {
		return createValid200Response(req), nil
	} else if mock.ProduceErrorResponse {
		return createError500Response(req), nil
	} else {
		return nil, errors.New("Always fails with no response")
	}
}

func createValid200Response(req *http.Request) *http.Response {
	body := "Hello world"
	t := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        make(http.Header, 0),
	}
	return t
}

func createError500Response(req *http.Request) *http.Response {
	var rserror RSpaceError = createTestRspace500Err()
	body, _ := json.Marshal(rserror)
	t := &http.Response{
		Status:        "500 Server Error",
		StatusCode:    500,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(string(body))),
		ContentLength: int64(len(string(body))),
		Request:       req,
		Header:        make(http.Header, 0),
	}
	return t
}
func createTestRspace500Err() RSpaceError {
	return RSpaceError{"Server error",
		500, 50001, "Error", make([]string, 0), "2020-05-03T12:34:56", 100}
}

func (mock *MockHttpClient) GetInvocationCount() int {
	return mock.numberOfInvocations
}

// retries with no client response
func TestRetryNoClientResponse(t *testing.T) {
	// it'll pass after 5 attempts
	mockClient := &MockHttpClient{PassOnNthInvocation: 5}
	assertIntEquals(t, 0, mockClient.GetInvocationCount(), "")

	retryClientEx, err := RetryClientExNew(0, mockClient)
	assertNotNil(t, err, "shouldn't be able to make a retry client with zero retries")

	retryClientEx, err = RetryClientExNew(3, mockClient)
	assertNotNil(t, retryClientEx, "should be created")

	resp, err := retryClientEx.Do(&http.Request{})
	assertNotNil(t, err, "Should always fail")
	assertIntEquals(t, 3, mockClient.GetInvocationCount(), "unexpected number of invocations")

	// try another 3 goes, should pass wno
	resp, err = retryClientEx.Do(&http.Request{})
	assertNotNil(t, resp, "resp should now be present")
	assertIntEquals(t, 5, mockClient.GetInvocationCount(), "unexpected number of invocations")
}

func TestRetryWithClientResponse(t *testing.T) {

	mockClient := &MockHttpClient{PassOnNthInvocation: 5, ProduceErrorResponse: true}
	retryClientEx, _ := RetryClientExNew(3, mockClient)
	retryClientEx.Do(&http.Request{})
	assertIntEquals(t, 3, mockClient.GetInvocationCount(), "unexpected number of invocations")
}
