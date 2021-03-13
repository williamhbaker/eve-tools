package lib

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

const namesFragment = "https://esi.evetech.net/v3/universe/names"

// AddNames adds names to a list of items, which presumably do not yet have names
func (e *Esi) AddNames(items map[int]*models.OrderItem, pageSize int) {
	list := []int{}

	for itemID := range items {
		list = append(list, itemID)

		if len(list) == pageSize {
			e.getNamesFromIDList(list, items)
			list = []int{}
		}
	}

	e.getNamesFromIDList(list, items)
}

func (e *Esi) getNamesFromIDList(l []int, i map[int]*models.OrderItem) {
	names := e.ItemNameList(l)

	for id, name := range names {
		i[id].Name = name
	}
}

// ItemNameList gets a list of item names from integer item ids
func (e *Esi) ItemNameList(list []int) map[int]string {
	output := make(map[int]string)
	body := fmt.Sprintf("%v", uniqueInts(list))
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

func uniqueInts(i []int) []int {
	seen := make(map[int]struct{})
	output := []int{}

	for _, v := range i {
		_, ok := seen[v]
		if !ok {
			seen[v] = struct{}{}
			output = append(output, v)
		}
	}

	return output
}
