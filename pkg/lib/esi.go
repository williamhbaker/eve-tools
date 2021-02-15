package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const ordersFragment = "https://esi.evetech.net/v1/markets/%d/orders?page=%d"

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

	for {
		url := fmt.Sprintf(ordersFragment, stationID, idx)
		resBytes, status, _ := e.get(url)

		res := []map[string]interface{}{}

		json.Unmarshal(resBytes, &res)
		fmt.Println(res[0]["is_buy_order"])

		fmt.Printf("Doing orders for page %d ...\n", idx)

		if status == 404 {
			break
		}

		if idx > 5 {
			break
		}

		idx++
	}

	return []string{"ehhlo"}
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
