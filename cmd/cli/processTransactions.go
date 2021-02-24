package main

import "github.com/wbaker85/eve-tools/pkg/lib"

func (app *application) processTransactions() {
	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)

	aggs := lib.MakeAggregates(transactions)

	d := []*lib.Aggregate{}

	for _, val := range aggs {
		d = append(d, val)
	}

	lib.SaveJSON("./transaction_aggregations.json", d)
}
