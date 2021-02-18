package lib

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/wbaker85/eve-tools/pkg/models"
)

const ordersFragment = "https://esi.evetech.net/v1/markets/%d/orders?page=%d"

// AllOrders iterates through all orders in the provided region and returns
// a slice containing the item id, buy price, and sell price
func (e *Esi) AllOrders(regionID, pageLimit int) []map[string]interface{} {
	idx := 1
	resultList := []map[string]interface{}{}

	for {
		fmt.Printf("Getting orders for page %d...\n", idx)

		if pageLimit > 0 && idx > pageLimit {
			break
		}

		url := fmt.Sprintf(ordersFragment, regionID, idx)
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
func AggregateOrders(items []map[string]interface{}, station int) map[int]*models.OrderItem {
	output := make(map[int]*models.OrderItem)

	for _, i := range items {
		itemID := int(i["type_id"].(float64))
		isBuy := i["is_buy_order"].(bool)
		price := i["price"].(float64)
		stationID := int(i["location_id"].(float64))

		if stationID != station {
			continue
		}

		_, ok := output[itemID]
		if !ok {
			output[itemID] = &models.OrderItem{
				ID:        itemID,
				BuyPrice:  0,
				SellPrice: 0,
			}
		}

		thisItem := output[itemID]

		if isBuy {
			if thisItem.BuyPrice == 0 {
				thisItem.BuyPrice = price
			} else {
				thisItem.BuyPrice = math.Max(thisItem.BuyPrice, price)
			}
		} else {
			if thisItem.SellPrice == 0 {
				thisItem.SellPrice = price
			} else {
				thisItem.SellPrice = math.Min(thisItem.SellPrice, price)
			}
		}
	}

	return output
}
