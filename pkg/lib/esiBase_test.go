package lib

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	c := &TestClient{DoFunc: func(r *http.Request) (*http.Response, error) {
		agent := r.Header.Get("User-Agent")
		url := r.URL.String()

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(agent + " - " + url)),
		}, nil
	},
	}

	s := "user@addr.com"
	u := "https://www.whatever.com"

	e := esi{client: c, userAgentString: s}

	res, _, _ := e.get(u)
	got := string(res)
	want := s + " - " + u

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
