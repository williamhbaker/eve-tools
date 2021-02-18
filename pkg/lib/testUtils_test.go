package lib

import "net/http"

type testClient struct {
	doFunc func(r *http.Request) (*http.Response, error)
}

func (t *testClient) Do(r *http.Request) (*http.Response, error) {
	return t.doFunc(r)
}

func newTestClient(doFunc func(r *http.Request) (*http.Response, error)) *testClient {
	return &testClient{doFunc: doFunc}
}
