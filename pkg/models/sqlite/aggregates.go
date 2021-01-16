package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	// For working with sqlite
	_ "github.com/mattn/go-sqlite3"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// AggregateModel deals with aggregated transaction data
type AggregateModel struct {
	DB *sql.DB
}

// LoadData re-initializes the database table of transaction aggregates and populates
// it with the provided map of aggregated transactions
func (a *AggregateModel) LoadData(d map[string]*models.Aggregate) {
	a.init()

	s := []*models.Aggregate{}

	for _, val := range d {
		s = append(s, val)
	}

	a.addMany(s)
}

// GetData returns slice of aggregated data - each aggregate is from the relevant
// transactions.
func (a *AggregateModel) GetData() []*models.Aggregate {
	stmt := `SELECT * FROM aggregates`

	output := []*models.Aggregate{}

	rows, err := a.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		a := &models.Aggregate{}
		err = rows.Scan(&a.Name, &a.Bought, &a.Sold, &a.Tax, &a.Spent, &a.Earned, &a.Profit)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, a)
	}

	return output
}

func (a *AggregateModel) addMany(aggs []*models.Aggregate) {
	stmt := `INSERT INTO aggregates (name, bought, sold, tax, spent, earned, profit) VALUES `

	for _, row := range aggs {
		sqlStr := "(%q, %d, %d, %f, %f, %f, %f),"
		stmt += fmt.Sprintf(
			sqlStr,
			row.Name,
			row.Bought,
			row.Sold,
			row.Tax,
			row.Spent,
			row.Earned,
			row.Profit,
		)
	}

	stmt = stmt[0 : len(stmt)-1]

	_, err := a.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *AggregateModel) init() {
	stmt := `CREATE TABLE aggregates (
		name   varchar(50) PRIMARY KEY,
		bought int,
		sold   int,
		tax    FLOAT,
		spent  FLOAT,
		earned FLOAT,
		profit FLOAT
	);`

	exists := `SELECT name FROM sqlite_master WHERE type='table' AND name='aggregates';`

	_, err := a.DB.Exec(exists)

	if err == nil {
		a.DB.Exec(`DROP TABLE aggregates`)
	}

	a.DB.Exec(stmt)
}
