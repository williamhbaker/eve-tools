package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

// ClientSecretModel is for interacting with client secrets in the database.
type ClientSecretModel struct {
	DB *sql.DB
}

// RegisterSecret replaces the current client ID in the database with the new one
func (c *ClientSecretModel) RegisterSecret(id string) {
	c.DB.Exec(`DROP TABLE client_secret`)

	_, err := c.DB.Exec(`CREATE TABLE client_secret (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			value VARCHAR(50)
		)`)

	if err != nil {
		log.Fatal(err)
	}

	stmt := `INSERT INTO client_secret (value) VALUES (%q)`

	_, err = c.DB.Exec(fmt.Sprintf(stmt, id))
	if err != nil {
		log.Fatal(err)
	}
}

// GetSecret returns the client ID from the database. There can be only one.
func (c *ClientSecretModel) GetSecret() string {
	var output string

	stmt := `SELECT value FROM client_secret`

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
