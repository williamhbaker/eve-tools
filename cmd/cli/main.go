package main

import (
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

const forgeRegionID = 10000002
const jitaStationID = 60003760
const perimiterTTTStationID = 1028858195912

type app struct {
	transactions *sqlite.TransactionModel
	parser       interface {
		ParseTransactions() []*models.Transaction
	}
}

func main() {
	// db, _ := sql.Open("sqlite3", "./data.db")
	// defer db.Close()

	// file, _ := os.Open("./transaction_export.csv")
	// defer file.Close()

	// app := app{
	// 	transactions: &sqlite.TransactionModel{DB: db},
	// 	parser:       &csvparser.TransactionParser{File: file},
	// }

	// transactions := app.parser.ParseTransactions()

	// app.transactions.LoadData(transactions)

	// aggs := lib.MakeAggregates(transactions)

	// d := []*lib.Aggregate{}

	// for _, val := range aggs {
	// 	d = append(d, val)
	// }

	// lib.SaveJSON("./transaction_aggregations.json", d)

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: "wbaker@gmail.com",
	}

	orders := api.AllOrders(forgeRegionID, 2)
	aggregates := lib.AggregateOrders(orders, jitaStationID)

	api.AddNames(aggregates)
}
