package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// AggregateModel deals with aggregated transaction data
type AggregateModel struct {
	DB *sql.DB
}

// AddMany adds many aggregates to the database
func (a *AggregateModel) AddMany(aggs []*models.Aggregate) {
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

	fmt.Println(stmt)

	_, err := a.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

}
