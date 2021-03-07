package main

import (
	"fmt"
	"net/url"
)

func loginURL(callbackURL, clientID, scopes string) string {
	loginURL := fmt.Sprintf("https://login.eveonline.com/v2/oauth/authorize/"+
		"?response_type=code"+
		"&redirect_uri=%v"+
		"&client_id=%v"+
		"&scope=%v"+
		"&state=%v",
		callbackURL, clientID, url.QueryEscape(scopes), "someState")

	return loginURL
}
