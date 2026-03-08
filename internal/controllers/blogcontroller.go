package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
		INSERT INTO posts (title, cover, author, fullname, tags)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id

	`

	err = tx.QueryRow(context.Background(), query, post.Title, post.Cover, post.Author, post.FullName, strings.Join(post.Tags, ", ")).Scan(&post.ID)
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

	pagination := "page:" + offset

	// get pagination from redis
	cached, err := db.RedisClient.Get(db.Ctx, pagination).Result()
	if err == nil {
		var posts []models.Post
		err = json.Unmarshal([]byte(cached), &posts)

		if err == nil {
			return posts, err
		}
	}

	query := fmt.Sprintf(`
		SELECT
			p.id, p.title, p.cover, p.author, p.fullname, p.date, p.visible, p.tags,
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

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var tags string
		var section models.BlogPost
		var sectionPosition *int
		var sectionType *string
		var sectionContent *string

		err = rows.Scan(
			&post.ID, &post.Title, &post.Cover, &post.Author, &post.FullName, &post.Date, &post.Visible, &tags,
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

	// save pagination if not exist in redis
	jsonData, _ := json.Marshal(posts)
	db.RedisClient.Set(db.Ctx, pagination, string(jsonData), time.Minute*2)

	return posts, nil
}

func GetOnePostFromDB(uuid string) (models.Post, error) {

	cacheKey := "post:" + uuid

	// check in redis first
	cached, err := db.RedisClient.Get(db.Ctx, cacheKey).Result()
	if err == nil {
		// return cached post -> from redis to front
		var post models.Post
		err = json.Unmarshal([]byte(cached), &post)

		// increment views in redis
		post.Views += 1
		jsonData, _ := json.Marshal(post)
		db.RedisClient.Set(db.Ctx, cacheKey, jsonData, time.Hour)

		if err == nil {
			return post, nil
		}
	}

	query := `
		SELECT id, title, cover, date, visible, tags, author, fullname, views FROM posts WHERE id = $1
	`
	var post models.Post
	var tags string

	err = db.Conn.QueryRow(context.Background(), query, uuid).Scan(&post.ID, &post.Title, &post.Cover, &post.Date, &post.Visible, &tags, &post.Author, &post.FullName, &post.Views)

	if err != nil {
		return post, err
	}

	post.Tags = append(post.Tags, strings.Split(tags, ", ")...)

	// find sections

	querySections := `
			SELECT id, position, type, content FROM blog_posts WHERE post_id = $1 ORDER BY position ASC
		`

	rows, err := db.Conn.Query(context.Background(), querySections, uuid)
	if err != nil {
		return models.Post{}, err
	}

	for rows.Next() {
		var section models.BlogPost

		if err = rows.Scan(&section.ID, &section.Position, &section.Type, &section.Content); err != nil {
			return models.Post{}, err
		}

		post.Sections = append(post.Sections, section)
	}

	// increment views
	post.Views += 1
	_, err = db.Conn.Exec(context.Background(), "UPDATE posts SET views = $1 WHERE ID = $2", post.Views, post.ID)

	// if post don't exist in redis, save it
	jsonData, _ := json.Marshal(post)
	db.RedisClient.Set(db.Ctx, cacheKey, jsonData, time.Hour)

	return post, nil
}

func GetPostsByTags(tag string, offset string) ([]models.Post, error) {

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
		LIMIT 5 OFFSET $2
		`

	rows, err := db.Conn.Query(context.Background(), query, "%"+tag+"%", offset)
	if err != nil {
		return nil, err
	}

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

func GetTopPosts(limit int) ([]models.Post, error) {
	cackeKey := fmt.Sprintf("top_posts:%d", limit)

	// check in redis first
	cached, err := db.RedisClient.Get(db.Ctx, cackeKey).Result()
	if err == nil {
		var posts []models.Post
		err = json.Unmarshal([]byte(cached), &posts)

		if err == nil {
			return posts, nil
		}
	}

	query := `
		SELECT id, title, cover, fullname, tags, date
		FROM posts
		ORDER BY views DESC
		LIMIT $1
	`

	rows, err := db.Conn.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var tags string

		err = rows.Scan(&post.ID, &post.Title, &post.Cover, &post.FullName, &tags, &post.Date)
		if err != nil {
			return nil, err
		}

		post.Tags = strings.Split(tags, ", ")

		posts = append(posts, post)
	}

	// save in redis
	jsonData, _ := json.Marshal(posts)
	db.RedisClient.Set(db.Ctx, cackeKey, string(jsonData), time.Hour*10)

	return posts, nil
}

func GetPostsByAuthor(author string, offset string) ([]models.Post, error) {
	query := `
		SELECT id, title, cover, date, visible, tags
		FROM posts
		WHERE author = $1
		ORDER BY date DESC
		LIMIT 5 OFFSET $2
	`

	rows, err := db.Conn.Query(context.Background(), query, author, offset)
	if err != nil {
		return nil, err
	}

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		var tags string

		err = rows.Scan(&post.ID, &post.Title, &post.Cover, &post.Date, &post.Visible, &tags)
		if err != nil {
			return nil, err
		}

		post.Tags = strings.Split(tags, ", ")

		posts = append(posts, post)
	}

	return posts, nil
}

func UpdatePost(post models.Post, user string) (bool, error) {
	cacheKey := "post:" + post.ID

	query := `
		UPDATE posts
		SET
			title = $1,
			cover = $2,
			tags = $3
		WHERE
			id = $4
		AND
			author = $5

	`

	_, err := db.Conn.Exec(context.Background(),
		query,
		post.Title,
		post.Cover,
		strings.Join(post.Tags, ", "),
		post.ID,
		user,
	)
	if err != nil {
		return false, err
	}

	/*
	 * Update post on redis cache
	 */
	db.RedisClient.Del(db.Ctx, cacheKey)

	return true, nil
}
