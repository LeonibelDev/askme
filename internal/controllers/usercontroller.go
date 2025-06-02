package controllers

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func GetUser(email string) (models.DBUser, error) {

	db.DataBaseConn()

	query := `
		SELECT * FROM users WHERE email = $1
	`
	var user models.DBUser

	rows, err := db.Conn.Query(context.Background(), query, email)
	if err != nil {
		return user, err
	}

	if !rows.Next() {
		return user, errors.New("user not found")
	}

	// social
	var twitter *string
	var github *string
	var instagram *string

	err = rows.Scan(&user.Id, &user.Fullname, &user.Username, &user.Email, &user.Password, &user.Role, &user.Resume, &user.Is_verified, &user.Created_at, &twitter, &github, &instagram, &user.External_link)
	if err != nil {
		return user, err
	}

	if twitter != nil {
		user.Twitter = *twitter
	}
	if github != nil {
		user.Github = *github
	}
	if instagram != nil {
		user.Instagram = *instagram
	}

	return user, nil
}

func UserExist(email string) bool {

	db.DataBaseConn()

	query := `
		SELECT id FROM users 
		WHERE email = $1 LIMIT 1
	`

	var id int
	err := db.Conn.QueryRow(context.Background(), query, email).Scan(&id)

	return bool(err == nil)
}

func CreateUser(user models.DBUser) bool {

	db.DataBaseConn()

	query := `
		INSERT INTO users (fullname, username, email, password) 
		VALUES ($1, $2, $3, $4)
	`

	err := db.Conn.QueryRow(context.Background(), query, user.Fullname, user.Username, user.Email, user.Password).Scan()

	return err == nil || err == pgx.ErrNoRows
}
