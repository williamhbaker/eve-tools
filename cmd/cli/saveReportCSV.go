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
			"group",
			"category",
			"meta",
			"tech",
			"sell_price",
			"buy_price",
			"margin",
			"recent_orders",
			"orders_avg",
			"volume_avg",
			"max_profit",
			"num_days",
			"year_max_sell",
			"year_min_sell",
			"year_max_buy",
			"year_min_buy",
			"relative_buy",
			"relative_sell",
		},
	}

	for _, item := range data {
		thisRecord := []string{
			fmt.Sprintf("%s", item.name),
			fmt.Sprintf("%d", item.itemID),
			fmt.Sprintf("%s", item.group),
			fmt.Sprintf("%s", item.category),
			fmt.Sprintf("%d", item.meta),
			fmt.Sprintf("%s", item.tech),
			fmt.Sprintf("%.2f", item.sellPrice),
			fmt.Sprintf("%.2f", item.buyPrice),
			fmt.Sprintf("%.2f", item.margin),
			fmt.Sprintf("%d", item.recentOrders),
			fmt.Sprintf("%d", item.ordersAvg),
			fmt.Sprintf("%d", item.volumeAvg),
			fmt.Sprintf("%.2f", item.maxProfit),
			fmt.Sprintf("%d", item.numDays),
			fmt.Sprintf("%.2f", item.yearMaxSell),
			fmt.Sprintf("%.2f", item.yearMinSell),
			fmt.Sprintf("%.2f", item.yearMaxBuy),
			fmt.Sprintf("%.2f", item.yearMinBuy),
			fmt.Sprintf("%.2f", item.relativeBuy),
			fmt.Sprintf("%.2f", item.relativeSell),
		}

		records = append(records, thisRecord)
	}

	lib.SaveCSV(path, records)
}
