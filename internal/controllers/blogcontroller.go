package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func SavePost(post models.Post) (bool, error) {

	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return false, err
	}

	// rollback transaction if error
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	// First insert data in post table
	query := `
		INSERT INTO posts (title, cover, author, tags) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id

	`

	err = tx.QueryRow(context.Background(), query, post.Title, post.Cover, post.Author, strings.Join(post.Tags, ", ")).Scan(&post.ID)
	if err != nil {
		return false, err
	}

	fmt.Println("Post ID: ", post.ID)

	// Then insert sections in po
	querySections := `
		INSERT INTO blog_posts (position, type, content, post_id) 
		VALUES ($1, $2, $3, $4)
	`

	for _, section := range post.Sections {
		_, err := tx.Exec(context.Background(), querySections, section.Position, section.Type, section.Content, post.ID)
		if err != nil {
			return false, fmt.Errorf("error inserting section: %v, error: %w", section, err)
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return false, err
	}

	defer tx.Conn().Close(context.Background())
	
	return true, nil
}

func GetAllPostsFromDB() ([]models.Post, error) {

	db.DataBaseConn()

	query := `
		SELECT * FROM posts
	`

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		// temporal variable
		var tags string

		if err = rows.Scan(&post.ID, &post.Title, &post.Cover, &post.Author, &post.Date, &post.Visible, &tags); err != nil {
			return nil, err
		}

		post.Tags = append(post.Tags, strings.Split(tags, ", ")...)

		posts = append(posts, post)
	}

	return posts, nil

}

func GetOnePostFromDB(uuid string) (models.Post, error) {

	db.DataBaseConn()

	query := `
		SELECT id, title, cover, date, visible, tags, author FROM posts WHERE id = $1
	`
	var post models.Post
	var tags string

	err := db.Conn.QueryRow(context.Background(), query, uuid).Scan(&post.ID, &post.Title, &post.Cover, &post.Date, &post.Visible, &tags, &post.Author)

	if err != nil {
		return models.Post{}, nil
	}

	post.Tags = append(post.Tags, strings.Split(tags, ", ")...)

	// find sections

	query_sections := `
			SELECT position, type, content FROM blog_posts WHERE post_id = $1 ORDER BY position ASC
		`

	rows, err := db.Conn.Query(context.Background(), query_sections, uuid)
	if err != nil {
		return models.Post{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var section models.BlogPost

		if err = rows.Scan(&section.Position, &section.Type, &section.Content); err != nil {
			return models.Post{}, err
		}

		post.Sections = append(post.Sections, section)
	}

	return post, nil
}
