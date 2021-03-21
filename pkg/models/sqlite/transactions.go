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

// LoadData loads the slice of transactions into the database. It clears the
// database first.
func (t *TransactionModel) LoadData(transactions []*models.Transaction) {
	t.init()
	t.addMany(transactions)
}

// GetAll returns every transaction from the database and returns a slice
func (t *TransactionModel) GetAll() []*models.Transaction {
	stmt := `SELECT * FROM transactions`

	output := []*models.Transaction{}

	rows, err := t.DB.Query(stmt)
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

func (t *TransactionModel) addMany(transactions []*models.Transaction) {
	if len(transactions) == 0 {
		return
	}

	var b strings.Builder
	stmt := `INSERT INTO transactions (date, name, quantity, price, tax, value, owner, station, region, client, type) VALUES `
	b.WriteString(stmt)

	for _, row := range transactions {
		sqlStr := "(%q, %q, %d, %f, %f, %f, %q, %q, %q, %q, %q),"
		b.WriteString(fmt.Sprintf(
			sqlStr,
			row.Date,
			row.Name,
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

// Profitable gets a list of items that have made the most profit, basically
func (t *TransactionModel) Profitable() []models.ProfitItem {
	stmt := `
SELECT	transactions.name,
	SUM(transactions.quantity) AS sold_qty,
	averages.buy_qty AS buy_qty,
	SUM(transactions.value) / SUM(transactions.quantity) AS avg_sell,
	averages.avg_buy,
	SUM(transactions.value) AS sold_val,
	averages.avg_buy*SUM(transactions.quantity) AS buy_val,
	SUM(transactions.value) - averages.avg_buy*SUM(transactions.quantity) AS profit
FROM transactions
INNER JOIN (
	SELECT name, ABS(SUM(value)/SUM(quantity)) as avg_buy, SUM(quantity) as buy_qty FROM transactions
	WHERE type = "Buy"
	GROUP BY name
) AS averages ON transactions.name = averages.name
WHERE transactions.type = "Sell"
GROUP BY transactions.name
ORDER BY profit DESC;
`

	output := []models.ProfitItem{}

	rows, err := t.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		t := models.ProfitItem{}
		err = rows.Scan(
			&t.Name,
			&t.SoldQty,
			&t.BuyQty,
			&t.AvgSell,
			&t.AvgBuy,
			&t.SoldVal,
			&t.BoughtVal,
			&t.Profit,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, t)
	}

	return output
}

// DateSales gets the total sale volume by date
func (t *TransactionModel) DateSales() []models.DateSale {
	stmt := `
	SELECT substr(date, 0, 11) AS date, SUM(value) AS total_sales FROM transactions
	WHERE TYPE = "Sell"
	GROUP BY substr(date, 0, 11)
	ORDER BY substr(date, 0, 11) DESC;
`

	output := []models.DateSale{}

	rows, err := t.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		t := models.DateSale{}
		err = rows.Scan(
			&t.Date,
			&t.TotalSales,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, t)
	}

	return output
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

	drop := `DROP TABLE transactions`

	t.DB.Exec(drop)

	_, err := t.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

}
