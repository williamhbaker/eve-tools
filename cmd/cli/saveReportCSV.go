package main

import (
	"fmt"

	"github.com/wbaker85/eve-tools/pkg/lib"
)

func saveReportCSV(path string, data []tradeItem) {
	records := [][]string{
		{
			"name",
			"item_id",
			"sell_price",
			"buy_price",
			"margin",
			"recent_orders",
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
			fmt.Sprintf("%d", item.recentOrders),
			fmt.Sprintf("%d", item.ordersAvg),
			fmt.Sprintf("%d", item.volumeAvg),
			fmt.Sprintf("%.2f", item.maxProfit),
			fmt.Sprintf("%d", item.numDays),
		}

		records = append(records, thisRecord)
	}

	lib.SaveCSV(path, records)
}
