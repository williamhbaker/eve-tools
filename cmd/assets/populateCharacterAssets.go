package main

import (
	"encoding/json"
	"fmt"
)

type asset struct {
}

func (app *application) populateCharacterAssets(charID int) {
	var assets []map[string]interface{}
	d := app.authorizedRequest(fmt.Sprintf(assetsURL, charID), "GET")
	json.Unmarshal(d, &assets)

	for _, val := range assets {
		fmt.Printf("%#v\n", val)
	}
}
