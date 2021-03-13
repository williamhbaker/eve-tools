package main

import (
	"encoding/json"
	"fmt"
)

type asset struct {
}

func (app *application) populateCharacterAssets() {
	var assets []map[string]interface{}
	d := app.authorizedRequest(assetsURL, "GET", true)
	json.Unmarshal(d, &assets)

	for _, val := range assets {
		fmt.Printf("%#v\n", val)
	}
}
