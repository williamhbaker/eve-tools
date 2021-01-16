package main

import (
	"database/sql"
	"fmt"
	"os"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/csv"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

type app struct {
	aggregates   *sqlite.AggregateModel
	transactions *sqlite.TransactionModel
	parser       interface {
		ParseTransactions() []*models.Transaction
	}
}

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	file, _ := os.Open("./transaction_export.csv")
	defer file.Close()

	app := app{
		aggregates:   &sqlite.AggregateModel{DB: db},
		transactions: &sqlite.TransactionModel{DB: db},
		parser:       &csv.TransactionParser{File: file},
	}

	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)

	aggs := app.aggregates.Aggregate()

	d := []*models.Aggregate{}

	for _, val := range aggs {
		d = append(d, val)
	}

	sort.Slice(d, func(i, j int) bool {
		return d[i].Profit < d[j].Profit
	})

	for _, item := range d {
		fmt.Printf("%s - %.0f\n", item.Name, item.Profit)
	}
}
