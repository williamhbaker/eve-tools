package lib

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/wbaker85/eve-tools/pkg/models"
)

const ordersFragment = "https://esi.evetech.net/v1/markets/%d/orders?page=%d"

// type item struct {
// 	id        int
// 	sellPrice float64
// 	buyPrice  float64
// }

// All orders iterates through all orders at the provided station and returns
// a slice containing the item id, buy price, and sell price
func (e *esi) AllOrders(stationID int) []map[string]interface{} {
	idx := 1
	resultList := []map[string]interface{}{}

	for {
		url := fmt.Sprintf(ordersFragment, stationID, idx)
		resBytes, status, _ := e.get(url)

		res := []map[string]interface{}{}

		json.Unmarshal(resBytes, &res)
		resultList = append(resultList, res...)

		if status == 404 {
			break
		}

		idx++
	}

	return resultList
}

// AggregateOrders aggregates data for a list of active orders from a single station.
func AggregateOrders(items []map[string]interface{}, stationID int) map[int]*models.OrderItem {
	output := make(map[int]*models.OrderItem)

	for _, i := range items {
		itemID := int(i["type_id"].(float64))
		isBuy := i["is_buy_order"]
		price := i["price"].(float64)

		_, ok := output[itemID]
		if !ok {
			output[itemID] = &models.OrderItem{
				ID:        itemID,
				StationID: stationID,
				BuyPrice:  math.Inf(-1),
				SellPrice: math.Inf(1),
			}
		}

		thisItem := output[itemID]

		if isBuy == "true" {
			thisItem.BuyPrice = math.Max(thisItem.BuyPrice, price)
		} else {
			thisItem.SellPrice = math.Min(thisItem.SellPrice, price)
		}
	}

	return output
}
