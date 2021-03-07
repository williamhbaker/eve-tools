package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
)

const charIDURL = "https://login.eveonline.com/oauth/verify"

func (app *application) authorizedRequest(url, method string) string {
	t, refreshed := lib.CurrentToken(app.authToken.GetToken(), app.clientID.GetID(), app.clientSecret.GetSecret())
	if refreshed {
		app.authToken.RegisterToken(models.AuthToken{
			AccessToken:  t.AccessToken,
			ExpiresIn:    t.ExpiresIn,
			RefreshToken: t.RefreshToken,
			Issued:       t.Issued,
		})
	}

	if method == "GET" {
		return string(app.authorizedGet(url))
	}

	return ""
}

func (app *application) authorizedGet(u string) []byte {
	authString := fmt.Sprintf("Bearer %v", app.authToken.GetToken().AccessToken)

	reqURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {authString},
		},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data
}
