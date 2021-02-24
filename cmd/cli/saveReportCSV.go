package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func saveReportCSV(path string, data []tradeItem) {
	records := [][]string{
		{
			"name",
			"item_id",
			"sell_price",
			"buy_price",
			"margin",
			"ordersAvg",
			"volumeAvg",
			"maxProfit",
			"numDays",
		},
	}

	for _, item := range data {
		thisRecord := []string{
			fmt.Sprintf("%s", item.name),
			fmt.Sprintf("%d", item.itemID),
			fmt.Sprintf("%.2f", item.sellPrice),
			fmt.Sprintf("%.2f", item.buyPrice),
			fmt.Sprintf("%.2f", item.margin),
			fmt.Sprintf("%d", item.ordersAvg),
			fmt.Sprintf("%d", item.volumeAvg),
			fmt.Sprintf("%.2f", item.maxProfit),
			fmt.Sprintf("%d", item.numDays),
		}

		records = append(records, thisRecord)
	}

	file, _ := os.Create(path)

	w := csv.NewWriter(file)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
