package csv

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// TransactionModel contains the methods for parsing a jEveAsssets export CSV
// file containing a transaction history
type TransactionModel struct {
	File *os.File
}

func (t *TransactionModel) AggregateTransactions() map[string]*models.Aggregate {
	r := csv.NewReader(t.File)
	r.Read()

	items := make(map[string]*models.Aggregate)

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		var i *models.Aggregate

		if items[record[1]] == nil {
			items[record[1]] = &models.Aggregate{Name: record[1]}
		}
		i = items[record[1]]

		qty, _ := strconv.Atoi(record[2])
		transVal, _ := strconv.ParseFloat(record[5], 64)
		tax, _ := strconv.ParseFloat(record[4], 64)

		thisType := record[14]
		if thisType == "Buy" {
			i.Bought += qty
			i.Spent += transVal
		} else {
			i.Sold += qty
			i.Earned += transVal
			i.Tax += tax
		}

		i.Profit = i.Earned + i.Spent + i.Tax
	}

	return items
}
