package lib

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

const namesFragment = "https://esi.evetech.net/v3/universe/names"

// AddNames adds names to a list of items, which presumably do not yet have names
func (e *Esi) AddNames(items map[int]*models.OrderItem) {
	// needs to look like: [11489, 17255]
	// can take 1000 at a time

	list := []int{}

	for itemID, itemData := range items {
		list = append(list, itemID)

		// go until we have a list that is 1000 items long
		if len(list) == 1000 {
			fmt.Println(list, len(list), itemData)
			list = []int{}
		}
	}

	fmt.Println(list, len(list))
}

func (e *Esi) itemNameList(list []int) map[int]string {
	output := make(map[int]string)
	body := fmt.Sprintf("%v", list)
	body = strings.ReplaceAll(body, " ", ",")

	res, _, _ := e.post(namesFragment, body)

	data := []map[string]interface{}{}
	json.Unmarshal(res, &data)

	for _, i := range data {
		id := int(i["id"].(float64))
		name := i["name"].(string)
		output[id] = name
	}

	return output
}
