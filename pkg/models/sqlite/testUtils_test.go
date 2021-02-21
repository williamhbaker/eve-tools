package sqlite

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", "./testDB.db")
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		db.Close()
		os.Remove("./testDB.db")
	}
}
