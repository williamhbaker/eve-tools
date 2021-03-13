package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
)

type order struct {
	Name       string
	Duration   int     `json:"duration"`
	IsBuy      bool    `json:"is_buy_order"`
	Issued     string  `json:"issued"`
	LocationID int     `json:"location_id"`
	MinVolume  int     `json:"min_volume"`
	OrderID    int     `json:"order_id"`
	Price      float64 `json:"price"`
	Range      string  `json:"range"`
	RegionID   int     `json:"region_id"`
	State      string  `json:"state"`
	TypeID     int     `json:"type_id"`
	VolRemain  int     `json:"volume_remain"`
	VolTotal   int     `json:"volume_total"`
}

func (app *application) getCharacterOrders() {
	var charData map[string]interface{}
	d := app.authorizedRequest(charIDURL, "GET")
	json.Unmarshal(d, &charData)

	charID := int(charData["CharacterID"].(float64))

	var orders []order
	d = app.authorizedRequest(fmt.Sprintf(ordersURL, charID), "GET")
	json.Unmarshal(d, &orders)

	charOrders := make([]*models.CharacterOrder, len(orders))

	for idx, o := range orders {
		thisOrder := models.CharacterOrder(o)
		charOrders[idx] = &thisOrder
	}

	nameCharacterOrders(charOrders)

	for _, val := range charOrders {
		fmt.Printf("%#v\n", val)
	}
}

func nameCharacterOrders(o []*models.CharacterOrder) {
	ids := []int{}

	for _, val := range o {
		ids = append(ids, val.TypeID)
	}

	esi := lib.Esi{
		Client: http.DefaultClient,
	}

	nameList := esi.ItemNameList(ids)

	for _, val := range o {
		id := val.TypeID
		val.Name = nameList[id]
	}
}
