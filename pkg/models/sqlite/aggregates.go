package sqlite

import (
	"database/sql"
	"log"

	// For working with sqlite
	_ "github.com/mattn/go-sqlite3"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// AggregateModel deals with aggregated transaction data
type AggregateModel struct {
	DB *sql.DB
}

// Aggregate aggregates the transactions in the database and returns a map
// of the result.
func (a *AggregateModel) Aggregate() map[string]*models.Aggregate {
	t := a.loadTransactions()

	output := make(map[string]*models.Aggregate)

	for _, trans := range t {
		var i *models.Aggregate

		if output[trans.Name] == nil {
			output[trans.Name] = &models.Aggregate{Name: trans.Name}
		}
		i = output[trans.Name]

		thisType := trans.Type
		if thisType == "Buy" {
			i.Bought += trans.Quantity
			i.Spent += trans.Value
		} else {
			i.Sold += trans.Quantity
			i.Earned += trans.Value
			i.Tax += trans.Tax
		}

		i.Profit = i.Earned + i.Spent + i.Tax
	}

	return output
}

func (a *AggregateModel) loadTransactions() []*models.Transaction {
	stmt := `SELECT * FROM transactions`

	output := []*models.Transaction{}

	rows, err := a.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		t := &models.Transaction{}
		err = rows.Scan(
			&t.ID,
			&t.Date,
			&t.Name,
			&t.Quantity,
			&t.Price,
			&t.Tax,
			&t.Value,
			&t.Owner,
			&t.Station,
			&t.Region,
			&t.Client,
			&t.Type,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, t)
	}

	return output
}
