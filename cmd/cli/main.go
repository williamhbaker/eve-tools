package main

import (
	"database/sql"
	"flag"
	"fmt"
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
	var updatePrices bool
	var updateVolumes bool
	var processTransactions bool
	var uaString string

	flag.BoolVar(&updatePrices, "prices", false, "Pass flag as true to update item prices from the ESI API")
	flag.BoolVar(&updateVolumes, "volumes", false, "Pass flag as true to update item volumes from the ESI API")
	flag.BoolVar(&processTransactions, "transactions", true, "Pass flag as true to process the jEveAssets exported transcations file located at ./transaction_export.csv")
	flag.StringVar(&uaString, "ua", "user@domain.com", "The string to use as the user agent for ESI API calls - usually an email address")
	flag.Parse()

	fmt.Println("updateVolumes: ", updateVolumes)
	fmt.Println("updatePrices: ", updateVolumes)
	fmt.Println("uaString: ", uaString)

	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	file, _ := os.Open("./transaction_export.csv")
	defer file.Close()

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: uaString,
	}

	app := application{
		transactions: &sqlite.TransactionModel{DB: db},
		orders:       &sqlite.OrderModel{DB: db},
		itemAverages: &sqlite.ItemAverageVolumeModel{DB: db},
		api:          &api,
		parser:       &csvparser.TransactionParser{File: file},
	}

	if updatePrices {
		app.updateOrdersByRegion(forgeRegionID, jitaStationID, perimiterTTTStationID)
	}

	margins := app.orders.GetAllMargins(jitaStationID, perimiterTTTStationID)

	if updateVolumes {
		app.updateItemVolumesByRegion(forgeRegionID, margins)
	}

	volumes := app.itemAverages.GetVolumesForRegion(forgeRegionID)

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
