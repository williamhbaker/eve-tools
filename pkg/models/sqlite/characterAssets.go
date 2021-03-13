package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// CharacterAssetModel holds the database and methods for interacting
// with a character order.
type CharacterAssetModel struct {
	DB *sql.DB
}

// LoadData loads the slice of transactions into the database. It clears the
// database first.
func (c *CharacterAssetModel) LoadData(assets []models.CharacterAsset) {
	c.init()
	c.addMany(assets)
}

// GetAll returns every transaction from the database and returns a slice
func (c *CharacterAssetModel) GetAll() []models.CharacterAsset {
	stmt := `SELECT name, type_id, location_flag, location_id, location_type, quantity FROM character_assets`

	output := []models.CharacterAsset{}

	rows, err := c.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		c := models.CharacterAsset{}
		err = rows.Scan(
			&c.Name,
			&c.TypeID,
			&c.LocationFlag,
			&c.LocationID,
			&c.LocationType,
			&c.Quantity,
		)
		if err != nil {
			log.Fatal(err)
		}

		output = append(output, c)
	}

	return output
}

func (c *CharacterAssetModel) addMany(assets []models.CharacterAsset) {
	if len(assets) == 0 {
		return
	}

	var b strings.Builder
	stmt := `INSERT INTO character_assets (name, type_id, location_flag, location_id, location_type, quantity) VALUES `
	b.WriteString(stmt)

	for _, row := range assets {
		sqlStr := "(%q, %d, %q, %d, %q, %d),"
		b.WriteString(fmt.Sprintf(
			sqlStr,
			row.Name,
			row.TypeID,
			row.LocationFlag,
			row.LocationID,
			row.LocationType,
			row.Quantity,
		))
	}

	stmt = b.String()
	stmt = stmt[0 : len(stmt)-1]

	_, err := c.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *CharacterAssetModel) init() {
	stmt := `CREATE TABLE character_assets (
		id						INTEGER PRIMARY KEY AUTOINCREMENT,
		name					VARCHAR(50),
		type_id				INT,
		location_flag	STRING,
		location_id		INT,
		location_type	STRING,
		quantity			INT
	);`

	drop := `DROP TABLE character_assets`

	c.DB.Exec(drop)

	_, err := c.DB.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}
