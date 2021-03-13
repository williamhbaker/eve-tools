package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// CharacterOrderModel holds the database and methods for interacting
// with a character order.
type CharacterOrderModel struct {
	DB *sql.DB
}

// LoadData loads the slice of transactions into the database. It clears the
// database first.
func (c *CharacterOrderModel) LoadData(orders []*models.CharacterOrder) {
	c.init()
	c.addMany(orders)
}

// GetAll returns every transaction from the database and returns a slice
func (c *CharacterOrderModel) GetAll() []*models.CharacterOrder {
	stmt := `SELECT name, duration, is_buy, issued, location_id, min_volume, order_id, price, range, region_id, state, type_id, volume_remaining, volume_total FROM character_orders`

	output := []*models.CharacterOrder{}

	rows, err := c.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		c := &models.CharacterOrder{}
		err = rows.Scan(
			&c.Name,
			&c.Duration,
			&c.IsBuy,
			&c.Issued,
			&c.LocationID,
			&c.MinVolume,
			&c.OrderID,
			&c.Price,
			&c.Range,
			&c.RegionID,
			&c.State,
			&c.TypeID,
			&c.VolRemain,
			&c.VolTotal,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, c)
	}

	return output
}

func (c *CharacterOrderModel) SellingInventory() []models.CharacterAsset {
	stmt := `SELECT type_id, name, SUM(volume_remaining) AS quantity
	FROM character_orders
	WHERE is_buy = 0
	GROUP BY name`

	output := []models.CharacterAsset{}

	rows, err := c.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		c := models.CharacterAsset{}
		err = rows.Scan(
			&c.TypeID,
			&c.Name,
			&c.Quantity,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, c)
	}

	return output
}

func (c *CharacterOrderModel) Orders(buy bool) []string {
	var buyFlag int

	if buy {
		buyFlag = 1
	} else {
		buyFlag = 0
	}

	stmt := `SELECT DISTINCT(name)
	FROM character_orders
	WHERE is_buy = %d AND volume_remaining > 0`

	output := []string{}

	rows, err := c.DB.Query(fmt.Sprintf(stmt, buyFlag))
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var c string
		err = rows.Scan(
			&c,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, c)
	}

	return output
}

func (c *CharacterOrderModel) addMany(orders []*models.CharacterOrder) {
	if len(orders) == 0 {
		return
	}

	var b strings.Builder
	stmt := `INSERT INTO character_orders (name, duration, is_buy, issued, location_id, min_volume, order_id, price, range, region_id, state, type_id, volume_remaining, volume_total) VALUES `
	b.WriteString(stmt)

	for _, row := range orders {
		sqlStr := "(%q, %d, %t, %q, %d, %d, %d, %f, %q, %d, %q, %d, %d, %d),"
		b.WriteString(fmt.Sprintf(
			sqlStr,
			row.Name,
			row.Duration,
			row.IsBuy,
			row.Issued,
			row.LocationID,
			row.MinVolume,
			row.OrderID,
			row.Price,
			row.Range,
			row.RegionID,
			row.State,
			row.TypeID,
			row.VolRemain,
			row.VolTotal,
		))
	}

	stmt = b.String()
	stmt = stmt[0 : len(stmt)-1]

	_, err := c.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CharacterOrderModel) init() {
	stmt := `CREATE TABLE character_orders (
		id								INTEGER PRIMARY KEY AUTOINCREMENT,
		name							VARCHAR(50),
		duration 					INT,
		is_buy 						BOOL,
		issued 						STRING,
		location_id				INT,
		min_volume				INT,
		order_id					INT,
		price							FLOAT,
		range							STRING,
		region_id					INT,
		state							STRING,
		type_id						INT,
		volume_remaining	INT,
		volume_total			INT
	);`

	drop := `DROP TABLE character_orders`

	c.DB.Exec(drop)

	_, err := c.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}

}
