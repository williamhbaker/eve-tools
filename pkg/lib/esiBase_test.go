package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	c := &testClient{doFunc: func(r *http.Request) (*http.Response, error) {
		agent := r.Header.Get("User-Agent")
		url := r.URL.String()

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(agent + " - " + url)),
		}, nil
	},
	}

	s := "user@addr.com"
	u := "https://www.whatever.com"

	e := Esi{Client: c, UserAgentString: s}

	res, _, _ := e.get(u)
	got := string(res)
	want := s + " - " + u

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestPost(t *testing.T) {
	c := &testClient{doFunc: func(r *http.Request) (*http.Response, error) {
		agent := r.Header.Get("User-Agent")
		url := r.URL.String()

		buf := bytes.Buffer{}
		buf.ReadFrom(r.Body)

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(agent + " - " + url + " - " + buf.String())),
		}, nil
	},
	}

	s := "user@addr.com"
	u := "https://www.whatever.com"
	b := "hello this is some test data"

	e := Esi{Client: c, UserAgentString: s}

	res, _, _ := e.post(u, b)
	got := string(res)
	want := s + " - " + u + " - " + b

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
