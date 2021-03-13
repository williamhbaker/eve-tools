package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
)

type itemRule struct {
	ItemName        string
	BuyTargetPrice  float64
	SellTargetPrice float64
	MaxInventory    int
	MinSellLotSize  int
}

func parseRules(path string) []itemRule {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	output := []itemRule{}

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		i := itemRule{}

		buyTargetPrice, _ := strconv.ParseFloat(strings.Replace(record[1], ",", "", -1), 64)
		sellTargetPrice, _ := strconv.ParseFloat(strings.Replace(record[2], ",", "", -1), 64)
		maxInventory, _ := strconv.Atoi(strings.Replace(record[3], ",", "", -1))
		minSellLotSize, _ := strconv.Atoi(strings.Replace(record[4], ",", "", -1))

		i.ItemName = record[0]
		i.BuyTargetPrice = buyTargetPrice
		i.SellTargetPrice = sellTargetPrice
		i.MaxInventory = maxInventory
		i.MinSellLotSize = minSellLotSize

		output = append(output, i)
	}

	return output
}
