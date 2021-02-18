package lib

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

type esi struct {
	client interface {
		Do(*http.Request) (*http.Response, error)
	}
	userAgentString string
}

func (e *esi) get(u string) ([]byte, int, error) {
	reqURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"User-Agent": {e.userAgentString},
		},
	}

	res, err := e.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data, res.StatusCode, nil
}
