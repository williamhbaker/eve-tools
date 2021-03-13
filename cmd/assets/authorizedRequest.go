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
const ordersURL = "https://esi.evetech.net/v1/characters/%d/orders"
const assetsURL = "https://esi.evetech.net/v5/characters/%d/assets/"

func (app *application) authorizedRequest(url, method string) []byte {
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
		return app.authorizedGet(url)
	}

	return []byte("")
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
