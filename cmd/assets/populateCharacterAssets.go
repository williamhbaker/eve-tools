package main

import (
	"encoding/json"

	"github.com/wbaker85/eve-tools/pkg/models"
)

type asset struct {
	Name         string
	TypeID       int    `json:"type_id"`
	LocationFlag string `json:"location_flag"`
	LocationID   int    `json:"location_id"`
	LocationType string `json:"location_type"`
	Quantity     int    `json:"quantity"`
}

func (a *asset) addName(n string) {
	a.Name = n
}

func (a *asset) id() int {
	return a.TypeID
}

func (app *application) populateCharacterAssets() {
	// Get the actual data from the ESI server
	var assets []asset
	d := app.authorizedRequest(assetsURL, "GET", true)
	json.Unmarshal(d, &assets)

	// Make a slice to hold interface representations of the assets
	nameableAssets := make([]nameable, len(assets))

	// Each element in the interface slice is a pointer to the actual asset
	for idx := range assets {
		nameableAssets[idx] = &assets[idx]
	}

	// Running this function calls each asset's addName method
	namer(nameableAssets)

	charAssets := make([]models.CharacterAsset, len(assets))

	for idx := range assets {
		charAssets[idx] = models.CharacterAsset(assets[idx])
	}

	app.characterAssets.LoadData(charAssets)
}
