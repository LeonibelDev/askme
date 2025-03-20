package controllers

import (
	"database/sql"

	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func SavePost(post models.Post) (bool, error) {

	db := db.Connection()
	defer db.Close()

	query := `
		INSERT INTO posts (title, cover, author, date, visible, tags, sections) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query, post.Title, post.Cover, post.Author, post.Date, post.Visible, post.Tags, post.Sections)

	if err == sql.ErrNoRows {
		return false, err
	}

	return true, nil
}
