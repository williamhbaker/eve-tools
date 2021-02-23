package lib

import (
	"net/http"
	"reflect"
	"testing"
)

type testClient struct {
	doFunc func(r *http.Request) (*http.Response, error)
}

func (t *testClient) Do(r *http.Request) (*http.Response, error) {
	return t.doFunc(r)
}

func newTestClient(doFunc func(r *http.Request) (*http.Response, error)) *testClient {
	return &testClient{doFunc: doFunc}
}

func assertInts(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertSlices(t *testing.T, got, want interface{}) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
