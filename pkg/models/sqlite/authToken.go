package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/wbaker85/eve-tools/pkg/models"
)

// AuthTokenModel is for interacting with auth tokens in the database.
type AuthTokenModel struct {
	DB *sql.DB
}

// RegisterToken replaces the current auth token in the database with the new one
func (a *AuthTokenModel) RegisterToken(t models.AuthToken) {
	a.DB.Exec(`DROP TABLE auth_token`)

	_, err := a.DB.Exec(`CREATE TABLE auth_token (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			access_token TEXT,
			expires_in INT,
			refresh_token TEXT,
			issued INT
		)`)

	if err != nil {
		log.Fatal(err)
	}

	stmt := `INSERT INTO auth_token (access_token, expires_in, refresh_token, issued) VALUES (%q, %d, %q, %d)`

	_, err = a.DB.Exec(fmt.Sprintf(stmt, t.AccessToken, t.ExpiresIn, t.RefreshToken, t.Issued))
	if err != nil {
		log.Fatal(err)
	}
}

// GetToken returns the auth token from the database. There can be only one.
func (a *AuthTokenModel) GetToken() models.AuthToken {
	var output models.AuthToken

	stmt := `SELECT access_token, expires_in, refresh_token, issued FROM auth_token`

	rows, err := a.DB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err = rows.Scan(&output.AccessToken, &output.ExpiresIn, &output.RefreshToken, &output.Issued)
		if err != nil {
			log.Fatal(err)
		}
	}

	return output
}
