package main

import (
	"fmt"

	"github.com/wbaker85/eve-tools/pkg/lib"
)

func saveTransactionsCSV(path string, data []*lib.Aggregate) {
	records := [][]string{
		{
			"name",
			"bought",
			"sold",
			"tax",
			"spent",
			"earned",
			"profit",
		},
	}

	for _, item := range data {
		thisRecord := []string{
			fmt.Sprintf("%s", item.Name),
			fmt.Sprintf("%d", item.Bought),
			fmt.Sprintf("%d", item.Sold),
			fmt.Sprintf("%.2f", item.Tax),
			fmt.Sprintf("%.2f", item.Spent),
			fmt.Sprintf("%.2f", item.Earned),
			fmt.Sprintf("%.2f", item.Profit),
		}

		records = append(records, thisRecord)
	}

	lib.SaveCSV(path, records)
}
