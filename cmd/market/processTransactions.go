package main

import (
	"sort"

	"github.com/wbaker85/eve-tools/pkg/lib"
)

func (app *application) processTransactions(path string) {
	transactions := app.parser.ParseTransactions()

	app.transactions.LoadData(transactions)

	aggs := lib.MakeAggregates(transactions)

	d := []*lib.Aggregate{}

	for _, val := range aggs {
		d = append(d, val)
	}

	sort.Slice(d, func(i, j int) bool {
		return d[i].Profit > d[j].Profit
	})

	saveTransactionsCSV(path, d)
}
