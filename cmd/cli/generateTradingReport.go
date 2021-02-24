package main

import (
	"sort"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func (app *application) generateTradingReport(reportPath string, margins []*models.MarginItem, volumes map[int]models.ItemAverageVolume) {
	output := []tradeItem{}

	for _, val := range margins {
		item := tradeItem{
			name:      val.Name,
			itemID:    val.ItemID,
			sellPrice: val.SellPrice,
			buyPrice:  val.BuyPrice,
			margin:    val.Margin,
		}

		item.ordersAvg = volumes[val.ItemID].OrdersAvg
		item.volumeAvg = volumes[val.ItemID].VolumeAvg
		item.numDays = volumes[val.ItemID].NumDays
		item.maxProfit = profitForItem(item)

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
