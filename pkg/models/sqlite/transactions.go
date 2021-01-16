package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// TransactionModel deals with transaction data
type TransactionModel struct {
	DB *sql.DB
}

// LoadData loads the slice of transactions into the database
func (t *TransactionModel) LoadData(transactions []*models.Transaction) {
	t.init()
	t.addMany(transactions)
}

func (t *TransactionModel) addMany(transactions []*models.Transaction) {
	var b strings.Builder
	stmt := `INSERT INTO transactions (date, name, quantity, price, tax, value, owner, station, region, client, type) VALUES `
	b.WriteString(stmt)

	for c, row := range transactions {
		if c%1000 == 0 {
			fmt.Printf("added item %d of %d\n", c, len(transactions))
		}

		sqlStr := "(%q, %q, %d, %f, %f, %f, %q, %q, %q, %q, %q),"
		b.WriteString(fmt.Sprintf(
			sqlStr,
			row.Name,
			row.Date,
			row.Quantity,
			row.Price,
			row.Tax,
			row.Value,
			row.Owner,
			row.Station,
			row.Region,
			row.Client,
			row.Type,
		))
	}

	stmt = b.String()
	stmt = stmt[0 : len(stmt)-1]

	_, err := t.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *TransactionModel) init() {
	stmt := `CREATE TABLE transactions (
		id				INTEGER PRIMARY KEY AUTOINCREMENT,
		name			VARCHAR(50),
		date			VARCHAR(20),
		quantity	INT,
		price			FLOAT,
		tax				FLOAT,
		value			FLOAT,
		owner			VARCHAR(50),
		station		VARCHAR(50),
		region		VARCHAR(50),
		client		VARCHAR(50),
		type			VARCHAR(5)
	);`

	exists := `SELECT name FROM sqlite_master WHERE type='table' AND name='transactions';`

	_, err := t.DB.Exec(exists)

	if err == nil {
		t.DB.Exec(`DROP TABLE transactions`)
	}

	_, err = t.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

}
