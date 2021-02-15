package lib

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type TestClient struct{}

func (t TestClient) Do(r *http.Request) (*http.Response, error) {
	agent := r.Header.Get("User-Agent")
	url := r.URL.String()

	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(agent + " - " + url)),
	}, nil
}

func TestGet(t *testing.T) {
	c := TestClient{}

	s := "wbaker@gmail.com"
	u := "http://www.google.com"

	e := esi{client: c, userAgentString: s}

	res, _, _ := e.get(u)
	got := string(res)
	want := s + " - " + u

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestAllOrders(t *testing.T) {
	c := http.DefaultClient

	s := "wbaker@gmail.com"
	u := "https://esi.evetech.net/v1/markets/10000002/orders?page=789"

	e := esi{client: c, userAgentString: s}

	_, status, _ := e.get(u)

	e.AllOrders(10000002)

	t.Errorf(strconv.Itoa(status))
}
