package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
)

const charIDURL = "https://login.eveonline.com/oauth/verify"
const ordersURL = "https://esi.evetech.net/v1/characters/%d/orders"
const assetsURL = "https://esi.evetech.net/v5/characters/%d/assets?page=%d"

func (app *application) authorizedRequest(url, method string, paginated bool) []byte {
	t, refreshed := lib.CurrentToken(app.authToken.GetToken(), app.clientID.GetID(), app.clientSecret.GetSecret())
	if refreshed {
		app.authToken.RegisterToken(models.AuthToken{
			AccessToken:  t.AccessToken,
			ExpiresIn:    t.ExpiresIn,
			RefreshToken: t.RefreshToken,
			Issued:       t.Issued,
		})
	}

	if method == "GET" && !paginated {
		if url == charIDURL {
			return app.authorizedGet(url)
		}
		return app.authorizedGet(fmt.Sprintf(url, app.charID))
	} else if method == "GET" && paginated {
		return app.paginatedAuthorizedGet(url)
	}

	return []byte("")
}

func (app *application) paginatedAuthorizedGet(u string) []byte {
	page := 1
	output := []byte{}

	for {
		url := fmt.Sprintf(u, app.charID, page)
		d := app.authorizedGet(url)

		if strings.Contains(string(d), "Requested page does not exist!") {
			return output
		}

		output = append(output, d...)
		page++
	}
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
