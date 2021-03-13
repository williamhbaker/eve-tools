package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

// ClientIDModel is for interacting with client ids in the database.
type ClientIDModel struct {
	DB *sql.DB
}

// RegisterID replaces the current client ID in the database with the new one
func (c *ClientIDModel) RegisterID(id string) {
	c.DB.Exec(`DROP TABLE client_id`)

	_, err := c.DB.Exec(`CREATE TABLE client_id (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			value VARCHAR(50)
		)`)

	if err != nil {
		log.Fatal(err)
	}

	stmt := `INSERT INTO client_id (value) VALUES (%q)`

	_, err = c.DB.Exec(fmt.Sprintf(stmt, id))
	if err != nil {
		log.Fatal(err)
	}
}

// GetID returns the client ID from the database. There can be only one.
func (c *ClientIDModel) GetID() string {
	var output string

	stmt := `SELECT value FROM client_id`

	rows, err := c.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&output)
		if err != nil {
			log.Fatal(err)
		}
	}

	return output
}
