package main

import (
	"database/sql"
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

	file, _ := os.Open("./transaction_export.csv")
	defer file.Close()

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: "wbaker@gmail.com",
	}

	app := application{
		transactions: &sqlite.TransactionModel{DB: db},
		orders:       &sqlite.OrderModel{DB: db},
		itemAverages: &sqlite.ItemAverageVolumeModel{DB: db},
		api:          &api,
		parser:       &csvparser.TransactionParser{File: file},
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

	app.processTransactions("./transactions.csv")
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
