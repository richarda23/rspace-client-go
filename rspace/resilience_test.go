package rspace

import (
	"errors"
	"net/http"
	"testing"
)

type MockHttpClient struct {
	numberOfInvocations int
	PassOnNthInvocation int
}

// MockClient counts  attethe number of times it is invoked and
// will pass on the  PassOnNthInvocation-th attempt
func (mock *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	mock.numberOfInvocations++
	if mock.numberOfInvocations >= mock.PassOnNthInvocation {
		return &http.Response{}, nil
	} else {
		return nil, errors.New("Always fails")
	}
}

func (mock *MockHttpClient) GetInvocationCount() int {
	return mock.numberOfInvocations
}

func TestRetryAlwaysFails(t *testing.T) {
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
