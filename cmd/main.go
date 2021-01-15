package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

// func main() {
// 	t := transactionAggregator{os.Open}
// 	items := t.aggregateTransactions("./transaction_export.csv")

// 	list := []itemData{}

// 	for _, val := range items {
// 		list = append(list, *val)
// 	}

// 	sort.Slice(list, func(i, j int) bool {
// 		return list[i].profit > list[j].profit
// 	})

// 	for _, val := range list {
// 		fmt.Printf("%s - %.0f\n", val.name, val.profit)
// 	}
// }

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()
	m := sqlite.AggregateModel{DB: db}

	t := transactionAggregator{os.Open}
	items := t.aggregateTransactions("./transaction_export.csv")

	iSlice := []*models.Aggregate{}

	for _, val := range items {
		iSlice = append(iSlice, val)
	}

	m.AddMany(iSlice[1:])

	fmt.Println(db)
}
