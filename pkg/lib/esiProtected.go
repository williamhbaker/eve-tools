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

	"github.com/wbaker85/eve-tools/pkg/models"
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

// CurrentToken evaluates the provided token. If its expired, a refreshed token
// is returned, along with a bool indicating if the token was refreshed or not
func CurrentToken(t models.AuthToken, clientID, clientSecret string) (Token, bool) {
	now := time.Now().Unix()

	if t.Issued+int64(t.ExpiresIn) < (now - 60) {
		return refreshToken(t, clientID, clientSecret), true
	}

	return Token{
		AccessToken:  t.AccessToken,
		ExpiresIn:    t.ExpiresIn,
		RefreshToken: t.RefreshToken,
		Issued:       t.Issued,
	}, false
}

func refreshToken(t models.AuthToken, clientID, clientSecret string) Token {
	bodyString := fmt.Sprintf("grant_type=refresh_token&refresh_token=%s", t.RefreshToken)
	authCredsString := fmt.Sprintf("%v:%v", clientID, clientSecret)
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

	var newToken Token
	err = json.Unmarshal(data, &newToken)
	if err != nil {
		log.Fatal(err)
	}

	newToken.Issued = time.Now().Unix()

	return newToken
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
