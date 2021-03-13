package sqlite

import "database/sql"

// CharacterOrderModel holds the database and methods for interacting
// with a character order.
type CharacterOrderModel struct {
	DB *sql.DB
}
