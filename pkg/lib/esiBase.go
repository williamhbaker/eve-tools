package lib

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Esi is a wrapper for actions involing interactions with the Eve Online ESI API.
type Esi struct {
	Client interface {
		Do(*http.Request) (*http.Response, error)
	}
	UserAgentString string
}

func (e *Esi) get(u string) ([]byte, int, error) {
	reqURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"User-Agent": {e.UserAgentString},
		},
	}

	res, err := e.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data, res.StatusCode, nil
}
