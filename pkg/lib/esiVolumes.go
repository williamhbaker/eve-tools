package lib

import (
	"encoding/json"
	"fmt"
)

const volumesFragment = `https://esi.evetech.net/v1/markets/%d/history?type_id=%d`

// VolumeForItem gets the volume information for a single item
func (e *Esi) VolumeForItem(regionID, itemID int) {
	u := fmt.Sprintf(volumesFragment, regionID, itemID)

	bytes, _, _ := e.get(u)

	var data []struct {
		Date       string  `json:"date"`
		Highest    float64 `json:"highest"`
		Lowest     float64 `json:"lowest"`
		OrderCount int     `json:"order_count"`
		Volume     int     `json:"volume"`
	}

	json.Unmarshal(bytes, &data)

	fmt.Println(data)
}
