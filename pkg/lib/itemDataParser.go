package lib

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

// ItemData represents the static data about an item.
type ItemData struct {
	TypeName string
	TypeID   int
	Group    string
	Category string
	Meta     int
	Tech     string
}

// ParseItemData parses an input file containing item data and returns a map
// relating the TypeIDs to the data
func ParseItemData(path string) map[int]ItemData {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	r.Read()

	output := make(map[int]ItemData)

	for {
		var i ItemData

		record, err := r.Read()
		if err != nil {
			break
		}

		name := record[0]
		id, _ := strconv.Atoi(strings.Replace(record[8], ",", "", -1))
		group := record[1]
		category := record[2]
		meta, _ := strconv.Atoi(record[5])
		tech := record[6]

		i.TypeName = name
		i.TypeID = id
		i.Group = group
		i.Category = category
		i.Meta = meta
		i.Tech = tech

		output[id] = i
	}

	return output
}
