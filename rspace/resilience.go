package rspace

import (
	"errors"
	"fmt"
	"net/http"
)

type TooManyRequests RSpaceError

// CientEx abstracts out the http.Client Do() method into an interface
// for decorating
type ClientEx interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryClientEx performs a fixed number of attempts to call clientEx.Do()
type RetryClientEx struct {
	retries int
	cli     ClientEx
}

func RetryClientExNew(retries int, wrappedRequest ClientEx) (*RetryClientEx, error) {
	if retries < 1 {
		return nil, errors.New("Number of retries must be >= 1")
	} else {
		return &RetryClientEx{retries, wrappedRequest}, nil
	}
}

func (ex RetryClientEx) Do(req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response
	for i := 0; i < ex.retries; i++ {
		resp, err = ex.cli.Do(req)
		if resp != nil {
			return resp, nil
		}
	}
	fmt.Println("error is " + err.Error())
	return nil, err
}
