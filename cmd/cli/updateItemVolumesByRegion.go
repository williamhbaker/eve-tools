package main

import "github.com/wbaker85/eve-tools/pkg/models"

func (app *application) updateItemVolumesByRegion(regionID int, items []*models.MarginItem) {
	itemIDs := []int{}
	for _, val := range items {
		itemIDs = append(itemIDs, val.ItemID)
	}

	volumes := app.api.VolumeForItems(regionID, itemIDs)

	app.itemAverages.LoadData(regionID, volumes)
}
