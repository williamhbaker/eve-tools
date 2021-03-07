package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const tokenURL = "https://login.eveonline.com/v2/oauth/token"
const tokenPostHost = "login.eveonline.com"

// Token is the token. Issued is the number of seconds since the unix epoch.
type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Issued       int64
}

// GetNewToken starts a server to listen for the callback from the ESI SSO.
// It returns the token when the auth flow is complete.
func GetNewToken(listenURL, clientID, clientSecret string) Token {
	c := make(chan Token)

	go func() {
		http.HandleFunc("/esi", func(w http.ResponseWriter, r *http.Request) {
			authCode := r.URL.Query()["code"][0]
			token := requestEsiToken(authCode, clientID, clientSecret)
			fmt.Fprintf(w, "Success, you can close this window now")

			c <- token
		})

		http.ListenAndServe(listenURL, nil)
	}()

	return <-c
}

func requestEsiToken(authCode, clientID, secret string) Token {
	bodyString := fmt.Sprintf("grant_type=authorization_code&code=%v", authCode)
	authCredsString := fmt.Sprintf("%v:%v", clientID, secret)
	encodedAuthCreds := "Basic " + base64.StdEncoding.EncodeToString([]byte(authCredsString))

	reqURL, _ := url.Parse(tokenURL)
	reqBody := ioutil.NopCloser(strings.NewReader(bodyString))
	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"Content-Type":  {"application/x-www-form-urlencoded"},
			"Authorization": {encodedAuthCreds},
			"Host":          {tokenPostHost},
		},
		Body: reqBody,
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var t Token
	err = json.Unmarshal(data, &t)
	if err != nil {
		log.Fatal(err)
	}

	t.Issued = time.Now().Unix()

	return t
}
