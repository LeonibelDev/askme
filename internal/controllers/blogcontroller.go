package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func SavePost(post models.Post) (string, error) {

	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return "", err
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
		return "", err
	}

	// Then insert sections in po
	querySections := `
		INSERT INTO blog_posts (position, type, content, post_id) 
		VALUES ($1, $2, $3, $4)
	`

	for _, section := range post.Sections {

		_, err := tx.Exec(context.Background(), querySections, section.Position, section.Type, section.Content, post.ID)

		if err != nil {
			return "", fmt.Errorf("error inserting section: %v, error: %w", section, err)
		}
	}

	if err = tx.Commit(context.Background()); err != nil {
		return "", err
	}

	return post.ID, nil
}

func GetAllPostsFromDB(offset string) ([]models.Post, error) {
	db.DataBaseConn()

	query := fmt.Sprintf(`
		SELECT 
			p.id, p.title, p.cover, p.author, p.date, p.visible, p.tags,
			b.position, b.type, b.content
		FROM posts p
			LEFT JOIN LATERAL (
				SELECT position, type, content
				FROM blog_posts
				WHERE post_id = p.id AND type = 'text'
				ORDER BY position ASC
				LIMIT 1
			) b ON true
			ORDER BY p.date DESC
			LIMIT 5 OFFSET %s
	`, offset)

	rows, err := db.Conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var tags string
		var section models.BlogPost
		var sectionPosition *int
		var sectionType *string
		var sectionContent *string

		err = rows.Scan(
			&post.ID, &post.Title, &post.Cover, &post.Author, &post.Date, &post.Visible, &tags,
			&sectionPosition, &sectionType, &sectionContent,
		)
		if err != nil {
			return nil, err
		}

		post.Tags = strings.Split(tags, ", ")

		// Check if section fields are not nil before adding the section to the post
		if sectionPosition != nil && sectionType != nil && sectionContent != nil {
			section.Position = *sectionPosition
			section.Type = *sectionType
			section.Content = *sectionContent

			post.Sections = append(post.Sections, section)
		} else {
			post.Sections = append(post.Sections, models.BlogPost{})
		}

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

func GetPostsByTags(tag string) ([]models.Post, error) {
	db.DataBaseConn()

	query := `
		SELECT 
			p.id, p.title, p.cover, p.author, p.date, p.visible, p.tags,
			b.position, b.type, b.content
		FROM posts p
		LEFT JOIN LATERAL (
			SELECT position, type, content
			FROM blog_posts
			WHERE post_id = p.id AND type = 'text'
			ORDER BY position ASC
			LIMIT 1
		) b ON true
		WHERE p.tags ILIKE $1
		ORDER BY p.date DESC
		`

	rows, err := db.Conn.Query(context.Background(), query, "%"+tag+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var tags string
		var section models.BlogPost
		var sectionPosition *int
		var sectionType *string
		var sectionContent *string

		err = rows.Scan(&post.ID, &post.Title, &post.Cover, &post.Author, &post.Date, &post.Visible, &tags,
			&sectionPosition, &sectionType, &sectionContent,
		)
		if err != nil {
			return nil, err
		}

		post.Tags = strings.Split(tags, ", ")

		if sectionPosition != nil && sectionType != nil && sectionContent != nil {
			section.Position = *sectionPosition
			section.Type = *sectionType
			section.Content = *sectionContent

			post.Sections = append(post.Sections, section)
		} else {
			post.Sections = append(post.Sections, models.BlogPost{})
		}

		posts = append(posts, post)
	}

	return posts, nil
}
