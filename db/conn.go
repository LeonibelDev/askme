package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/mattn/go-sqlite3"
)

var Conn *pgx.Conn

func DataBaseConn() error {
	var err error
	Conn, err = pgx.Connect(context.Background(), os.Getenv("PG_URI"))
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("Connected to database")
	return nil
}

func CreateTables() error {

	DataBaseConn()

	query :=
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			hashpassword TEXT NOT NULL,
			role TEXT DEFAULT 'admin',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			cover TEXT,
			author TEXT NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			visible BOOLEAN DEFAULT FALSE,
			tags TEXT,
			FOREIGN KEY (author) REFERENCES users(email) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS blog_posts (
			id SERIAL PRIMARY KEY,
			position INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL,
			post_id UUID,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS newsletter (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL UNIQUE,
			inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`

	_, err := Conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}
