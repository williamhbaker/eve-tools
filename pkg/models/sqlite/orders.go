package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func escapeString(s string) string {
	return strings.Replace(s, "\"", "\"\"", -1)
}

// OrderModel deals with database actions pertaining to orders.
type OrderModel struct {
	DB *sql.DB
}

// LoadData loads a slice of orders into the database for the given station ID.
// It clears the table for that station ID first - all order are stored in
// a table per station ID.
func (o *OrderModel) LoadData(stationID int, data map[int]*models.OrderItem) {
	o.init(stationID)
	o.addMany(stationID, data)
}

// GetAllMargins returns a list of margins for all of the items in the database
func (o *OrderModel) GetAllMargins(sellStationID, buyStationID int) []*models.MarginItem {
	output := []*models.MarginItem{}

	stmt := `
SELECT
	sell_station.item_id AS item_id,
	sell_station.name AS name,
	MAX(sell_station.buy_price, buy_station.buy_price) AS buy_price,
	sell_station.sell_price AS sell_price,
	(sell_station.sell_price - MAX(sell_station.buy_price, buy_station.buy_price)) / MAX(sell_station.buy_price, buy_station.buy_price) * 100 AS margin
FROM "%d_orders" AS sell_station
INNER JOIN "%d_orders" AS buy_station
	ON sell_station.item_id = buy_station.item_id
WHERE
	sell_station.sell_price > 0 AND
	MAX(sell_station.buy_price, buy_station.buy_price) > 0
ORDER BY margin DESC;`

	stmt = fmt.Sprintf(stmt, sellStationID, buyStationID)

	rows, err := o.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		i := &models.MarginItem{}
		err = rows.Scan(
			&i.ID,
			&i.Name,
			&i.BuyPrice,
			&i.SellPrice,
			&i.Margin,
		)

		output = append(output, i)
	}

	return output
}

func (o *OrderModel) addMany(stationID int, data map[int]*models.OrderItem) {
	if len(data) == 0 {
		return
	}

	var b strings.Builder
	stmt := fmt.Sprintf(`INSERT INTO "%d_orders" (item_id, name, sell_price, buy_price) VALUES `, stationID)
	b.WriteString(stmt)

	for _, item := range data {
		sqlStr := `(%d, "%s", %f, %f),`
		b.WriteString(fmt.Sprintf(sqlStr, item.ID, escapeString(item.Name), item.SellPrice, item.BuyPrice))
	}

	stmt = b.String()
	stmt = stmt[:len(stmt)-1]

	_, err := o.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (o *OrderModel) init(stationID int) {
	stmt := `CREATE TABLE "%d_orders" (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_id INT,
		name VARCHAR(50),
		sell_price FLOAT,
		buy_price FLOAT
	)`

	drop := `DROP TABLE "%d_orders"`

	o.DB.Exec(fmt.Sprintf(drop, stationID))

	_, err := o.DB.Exec(fmt.Sprintf(stmt, stationID))
	if err != nil {
		log.Fatal(err)
	}
}
