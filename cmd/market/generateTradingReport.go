package main

import (
	"sort"

	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
)

type tradeItem struct {
	name         string
	itemID       int
	group        string
	category     string
	meta         int
	tech         string
	sellPrice    float64
	buyPrice     float64
	margin       float64
	recentOrders int
	ordersAvg    int
	volumeAvg    int
	maxProfit    float64
	numDays      int
	yearMaxSell  float64
	yearMinSell  float64
	yearMaxBuy   float64
	yearMinBuy   float64
	relativeBuy  float64
	relativeSell float64
}

func (app *application) generateTradingReport(reportPath string, margins []*models.MarginItem, volumes map[int]models.ItemHistoryData, metaData map[int]lib.ItemData) {
	output := []tradeItem{}

	for _, val := range margins {
		item := tradeItem{
			name:         val.Name,
			itemID:       val.ItemID,
			sellPrice:    val.SellPrice,
			buyPrice:     val.BuyPrice,
			margin:       val.Margin,
			recentOrders: val.RecentOrders,
		}

		item.ordersAvg = volumes[val.ItemID].OrdersAvg
		item.volumeAvg = volumes[val.ItemID].VolumeAvg
		item.numDays = volumes[val.ItemID].NumDays
		item.yearMaxSell = volumes[val.ItemID].YearMaxSell
		item.yearMinSell = volumes[val.ItemID].YearMinSell
		item.yearMaxBuy = volumes[val.ItemID].YearMaxBuy
		item.yearMinBuy = volumes[val.ItemID].YearMinBuy
		item.relativeBuy = relativePrice(item.yearMaxBuy, item.yearMinBuy, item.buyPrice)
		item.relativeSell = relativePrice(item.yearMaxSell, item.yearMinSell, item.sellPrice)
		item.maxProfit = profitForItem(item)

		item.group = metaData[val.ItemID].Group
		item.category = metaData[val.ItemID].Category
		item.meta = metaData[val.ItemID].Meta
		item.tech = metaData[val.ItemID].Tech

		output = append(output, item)
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].maxProfit > output[j].maxProfit
	})

	saveReportCSV(reportPath, output)
}

func profitForItem(i tradeItem) float64 {
	return 0.5 * float64(i.volumeAvg) * (i.sellPrice - i.buyPrice)
}

func relativePrice(max, min, current float64) float64 {
	numerator := current - min
	denominator := max - min

	if denominator <= 0 || numerator <= 0 {
		return 0
	}

	return (numerator / denominator) * 100
}
