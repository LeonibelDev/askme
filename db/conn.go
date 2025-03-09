package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connection() *sql.DB {
	db, err := sql.Open("sqlite3", "db/users.db")
	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}

	return db
}

func VerifyIfDBExist() error {

	db := Connection()
	defer db.Close()

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			hashpassword TEXT,
			role TEXT DEFAULT "admin",
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)

	return err

}
