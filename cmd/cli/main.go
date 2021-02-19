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

type app struct {
	transactions *sqlite.TransactionModel
	orders       *sqlite.OrderModel
	parser       interface {
		ParseTransactions() []*models.Transaction
	}
}

func main() {
	// scrapeOrders()
	// processTransactions()

	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	app := app{
		transactions: &sqlite.TransactionModel{DB: db},
		orders:       &sqlite.OrderModel{DB: db},
	}

	orders := scrapeOrders(forgeRegionID, jitaStationID)

	app.orders.LoadData(jitaStationID, orders)
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

func scrapeOrders(regionID, stationID int) map[int]*models.OrderItem {
	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: "wbaker@gmail.com",
	}

	orders := api.AllOrders(regionID, -1)
	aggregates := lib.AggregateOrders(orders, stationID)

	api.AddNames(aggregates, 1000)

	return aggregates
}
