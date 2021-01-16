package csv

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// TransactionParser deals with parsing transactions from a file
type TransactionParser struct {
	File io.Reader
}

// ParseTransactions parses the transactions from a CSV file and returns the
// resulting slice.
func (t *TransactionParser) ParseTransactions() []*models.Transaction {
	r := csv.NewReader(t.File)
	r.Read()

	transactions := []*models.Transaction{}

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		i := &models.Transaction{}

		qty, _ := strconv.Atoi(record[2])
		price, _ := strconv.ParseFloat(record[3], 64)
		tax, _ := strconv.ParseFloat(record[4], 64)
		value, _ := strconv.ParseFloat(record[5], 64)

		i.Date = record[0]
		i.Name = record[1]
		i.Quantity = qty
		i.Price = price
		i.Tax = tax
		i.Value = value
		i.Owner = record[6]
		i.Station = record[7]
		i.Region = record[8]
		i.Client = record[13]
		i.Type = record[14]

		transactions = append(transactions, i)
	}

	return transactions
}
