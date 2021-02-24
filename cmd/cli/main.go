package main

import (
	"database/sql"
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

type app struct {
	transactions *sqlite.TransactionModel
	orders       *sqlite.OrderModel
	itemAverages *sqlite.ItemAverageVolumeModel
	parser       interface {
		ParseTransactions() []*models.Transaction
	}
}

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	// app := app{
	// 	transactions: &sqlite.TransactionModel{DB: db},
	// 	orders:       &sqlite.OrderModel{DB: db},
	// 	itemAverages: &sqlite.ItemAverageVolumeModel{DB: db},
	// }

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: "wbaker@gmail.com",
	}

	forgeOrders := api.AllOrders(forgeRegionID, 1)

	jitaPrices := lib.AggregateOrders(forgeOrders, jitaStationID)
	tttPrices := lib.AggregateOrders(forgeOrders, perimiterTTTStationID)

	api.AddNames(jitaPrices, 1000)
	api.AddNames(tttPrices, 1000)

	volumes := api.VolumeForItems(forgeRegionID, jitaPrices)

	for _, val := range volumes {
		fmt.Printf("%#v\n", val)
	}

	// app.orders.LoadData(jitaStationID, jitaPrices)
	// app.orders.LoadData(perimiterTTTStationID, tttPrices)

	// margins := app.orders.GetAllMargins(jitaStationID, perimiterTTTStationID)
	// lib.SaveJSON("./margins.json", margins)
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

// // "duration":90,"is_buy_order":false,"issued":"2021-02-15T20:28:58Z","location_id":60003760,"min_volume":1,"order_id":5926551566,"price":392.2,"range":"region","system_id":30000142,"type_id":2400,"volume_remain":349627,"volume_total":349627

// func saveCSV(path string, data []map[string]interface{}) {
// 	records := [][]string{
// 		{
// 			"duration",
// 			"is_buy_order",
// 			"issued",
// 			"location_id",
// 			"min_volume",
// 			"order_id",
// 			"price",
// 			"range",
// 			"system_id",
// 			"type_id",
// 			"volume_remain",
// 			"volume_total",
// 		},
// 	}

// 	for _, item := range data {
// 		thisRecord := []string{
// 			fmt.Sprintf("%v", item["duration"].(float64)),
// 			fmt.Sprintf("%v", item["is_buy_order"].(bool)),
// 			fmt.Sprintf("%v", item["issued"].(string)),
// 			fmt.Sprintf("%v", item["location_id"].(float64)),
// 			fmt.Sprintf("%v", item["min_volume"].(float64)),
// 			fmt.Sprintf("%v", item["order_id"].(float64)),
// 			fmt.Sprintf("%v", item["price"].(float64)),
// 			fmt.Sprintf("%v", item["range"].(string)),
// 			fmt.Sprintf("%v", item["system_id"].(float64)),
// 			fmt.Sprintf("%v", item["type_id"].(float64)),
// 			fmt.Sprintf("%v", item["volume_remain"].(float64)),
// 			fmt.Sprintf("%v", item["volume_total"].(float64)),
// 		}

// 		records = append(records, thisRecord)
// 	}

// 	file, _ := os.Create(path)

// 	w := csv.NewWriter(file)

// 	for _, record := range records {
// 		if err := w.Write(record); err != nil {
// 			log.Fatalln("error writing record to csv:", err)
// 		}
// 	}

// 	// Write any buffered data to the underlying writer (standard output).
// 	w.Flush()

// 	if err := w.Error(); err != nil {
// 		log.Fatal(err)
// 	}

// }
