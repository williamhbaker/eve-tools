package main

import (
	"database/sql"
	"os"

	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/csvparser"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

type application struct {
	transactions *sqlite.TransactionModel
	orders       *sqlite.OrderModel
	itemAverages *sqlite.ItemHistoryDataModel
	api          *lib.Esi
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

}

func (app *application) processTransactions() {
	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)
}
