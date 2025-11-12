package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/mattn/go-sqlite3"
)

var Conn *pgxpool.Pool

func DataBaseConn() error {
	var err error
	Conn, err = pgxpool.New(context.Background(), os.Getenv("PG_URI"))
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("Connected to database")
	return nil
}

func CreateTables() error {

	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			fullname TEXT NOT NULL,
			username TEXT NOT NULL UNIQUE,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role TEXT DEFAULT 'user',
			resume TEXT,
			is_verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			twitter TEXT,
			github TEXT,
			instagram TEXT,
			external_link TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			cover TEXT,
			author TEXT NOT NULL,
			fullname TEXT NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			visible BOOLEAN DEFAULT FALSE,
			tags TEXT,
			FOREIGN KEY (author) REFERENCES users(username) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS blog_posts (
			id SERIAL PRIMARY KEY,
			position INTEGER NOT NULL,
			type TEXT NOT NULL,
			content TEXT NOT NULL,
			post_id UUID,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS newsletter (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL UNIQUE,
			inserted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`}

	for _, q := range queries {
		_, err := Conn.Exec(context.Background(), q)
		if err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}

	return nil
}
