package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/lib"
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
	saveProfitCSV("./profit_list.csv", profitList)

	dateList := app.transactions.DateSales()
	saveDateCSV("./date_sales.csv", dateList)
}

func (app *application) processTransactions() {
	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)
}

func saveProfitCSV(path string, data []models.ProfitItem) {
	records := [][]string{
		{
			"name",
			"sold_qty",
			"buy_qty",
			"avg_sell",
			"avg_buy",
			"sold_val",
			"bought_val",
			"profit",
		},
	}

	for _, item := range data {
		thisRecord := []string{
			fmt.Sprintf("%s", item.Name),
			fmt.Sprintf("%d", item.SoldQty),
			fmt.Sprintf("%d", item.BuyQty),
			fmt.Sprintf("%.2f", item.AvgSell),
			fmt.Sprintf("%.2f", item.AvgBuy),
			fmt.Sprintf("%.2f", item.SoldVal),
			fmt.Sprintf("%.2f", item.BoughtVal),
			fmt.Sprintf("%.2f", item.Profit),
		}

		records = append(records, thisRecord)
	}

	lib.SaveCSV(path, records)
}

func saveDateCSV(path string, data []models.DateSale) {
	records := [][]string{
		{
			"date",
			"total_sales",
		},
	}

	for _, item := range data {
		thisRecord := []string{
			fmt.Sprintf("%s", item.Date),
			fmt.Sprintf("%.2f", item.TotalSales),
		}

		records = append(records, thisRecord)
	}

	lib.SaveCSV(path, records)
}
