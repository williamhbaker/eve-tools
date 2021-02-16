package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
)

const ordersFragment = "https://esi.evetech.net/v1/markets/%d/orders?page=%d"

type item struct {
	id        int
	sellPrice float64
	buyPrice  float64
}

type esi struct {
	client interface {
		Do(*http.Request) (*http.Response, error)
	}
	userAgentString string
}

// All orders iterates through all orders at the provided station and returns
// a slice containing the item id, buy price, and sell price
func (e *esi) AllOrders(stationID int) []string {
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

		if idx > 5 {
			break
		}

		idx++
	}

	aggregateOrders(resultList)

	return []string{"ehhlo"}
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

func (e *esi) get(u string) ([]byte, int, error) {
	reqURL, _ := url.Parse(u)

	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"User-Agent": {e.userAgentString},
		},
	}

	res, err := e.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data, res.StatusCode, nil
}
