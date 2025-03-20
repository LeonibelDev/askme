package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Connection() *sql.DB {
	db, err := sql.Open("sqlite3", "db/askme.db")
	if err != nil {
		panic(err)
	}

	if db == nil {
		panic("db nil")
	}

	return db
}

func CreateTables() error {

	db := Connection()
	defer db.Close()

	queris := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			hashpassword TEXT,
			role TEXT DEFAULT "admin",
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`,
		`
		CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			title TEXT,
			cover TEXT,
			author TEXT,
			date TIMESTAMP,
			visible BOOLEAN DEFAULT false,
			tags TEXT
		);
	`,
		`CREATE TABLE IF NOT EXISTS blog_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			position INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL,
			post_id INTEGER,
			FOREIGN KEY(post_id) REFERENCES posts(id)
		);
		`,
	}

	for _, query := range queris {
		db.Exec(query)
	}

	return nil
}
