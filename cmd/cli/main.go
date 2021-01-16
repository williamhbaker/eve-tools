package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models/csv"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

type app struct {
	aggregates   *sqlite.AggregateModel
	transactions *csv.TransactionModel
}

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	file, _ := os.Open("./transaction_export.csv")
	defer file.Close()

	app := app{
		aggregates:   &sqlite.AggregateModel{DB: db},
		transactions: &csv.TransactionModel{File: file},
	}

	items := app.transactions.AggregateTransactions()
	app.aggregates.LoadData(items)

	d := app.aggregates.GetData()

	sort.Slice(d, func(i, j int) bool {
		return d[i].Profit < d[j].Profit
	})

	for _, item := range d {
		fmt.Printf("%s - %.0f\n", item.Name, item.Profit)
	}
}
