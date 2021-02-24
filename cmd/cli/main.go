package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/csvparser"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

const forgeRegionID = 10000002
const jitaStationID = 60003760
const perimiterTTTStationID = 1028858195912

type application struct {
	transactions *sqlite.TransactionModel
	orders       *sqlite.OrderModel
	itemAverages *sqlite.ItemAverageVolumeModel
	api          *lib.Esi
	parser       interface {
		ParseTransactions() []*models.Transaction
	}
}

type tradeItem struct {
	name      string
	itemID    int
	sellPrice float64
	buyPrice  float64
	margin    float64
	ordersAvg int
	volumeAvg int
	maxProfit float64
	numDays   int
}

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: "wbaker@gmail.com",
	}

	app := application{
		transactions: &sqlite.TransactionModel{DB: db},
		orders:       &sqlite.OrderModel{DB: db},
		itemAverages: &sqlite.ItemAverageVolumeModel{DB: db},
		api:          &api,
	}

	// This makes a bunch of API calls - saves results to DB
	// app.updateOrdersByRegion(forgeRegionID, jitaStationID, perimiterTTTStationID)

	// Just database calls
	margins := app.orders.GetAllMargins(jitaStationID, perimiterTTTStationID)

	// API calls - lots of them! - also saves results to DB
	// app.updateItemVolumesByRegion(forgeRegionID, margins)

	// Just database calls
	volumes := app.itemAverages.GetVolumesForRegion(forgeRegionID)

	// This is what ends up saving a csv
	app.generateTradingReport("./report.csv", margins, volumes)

}

func (app *application) updateOrdersByRegion(regionID, sellStationID, buyStationID int) {
	orders := app.api.AllOrders(regionID, -1)

	sellStationPrices := lib.AggregateOrders(orders, sellStationID)
	buyStationPrices := lib.AggregateOrders(orders, buyStationID)

	app.api.AddNames(sellStationPrices, 1000)
	app.api.AddNames(buyStationPrices, 1000)

	app.orders.LoadData(sellStationID, sellStationPrices)
	app.orders.LoadData(buyStationID, buyStationPrices)
}

func (app *application) updateItemVolumesByRegion(regionID int, items []*models.MarginItem) {
	itemIDs := []int{}
	for _, val := range items {
		itemIDs = append(itemIDs, val.ItemID)
	}

	volumes := app.api.VolumeForItems(regionID, itemIDs)

	app.itemAverages.LoadData(regionID, volumes)
}

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

	saveReportCSV(reportPath, output)
}

func profitForItem(i tradeItem) float64 {
	return 0.5 * float64(i.volumeAvg) * (i.sellPrice - i.buyPrice)
}

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

func processTransactions() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	file, _ := os.Open("./transaction_export.csv")
	defer file.Close()

	app := app{
		transactions: &sqlite.TransactionModel{DB: db},
		parser:       &csvparser.TransactionParser{File: file},
	}

	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)

	aggs := lib.MakeAggregates(transactions)

	d := []*lib.Aggregate{}

	for _, val := range aggs {
		d = append(d, val)
	}

	lib.SaveJSON("./transaction_aggregations.json", d)
}
