package lib

import (
	"encoding/json"
	"fmt"
	"math"
)

const ordersFragment = "https://esi.evetech.net/v1/markets/%d/orders?page=%d"

type item struct {
	id        int
	sellPrice float64
	buyPrice  float64
}

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

func aggregateOrders(items []map[string]interface{}) map[int]*item {
	output := make(map[int]*item)

	for _, i := range items {
		itemID := int(i["type_id"].(float64))
		isBuy := i["is_buy_order"]
		price := i["price"].(float64)

		_, ok := output[itemID]
		if !ok {
			output[itemID] = &item{
				id:        itemID,
				buyPrice:  math.Inf(-1),
				sellPrice: math.Inf(1),
			}
		}

		thisItem := output[itemID]

		if isBuy == "true" {
			thisItem.buyPrice = math.Max(thisItem.buyPrice, price)
		} else {
			thisItem.sellPrice = math.Min(thisItem.sellPrice, price)
		}
	}

	return output
}
