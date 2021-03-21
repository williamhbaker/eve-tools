package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/csvparser"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

type application struct {
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

	app := application{
		transactions: &sqlite.TransactionModel{DB: db},
		parser:       &csvparser.TransactionParser{File: file},
	}

	app.processTransactions()

	profitList := app.transactions.Profitable()

	for _, val := range profitList {
		fmt.Println(val)
	}

	dateList := app.transactions.DateSales()

	for _, val := range dateList {
		fmt.Println(val)
	}

}

func (app *application) processTransactions() {
	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)
}
