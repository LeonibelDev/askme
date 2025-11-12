package controllers

import (
	"context"
	"errors"

	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func GetUser(email string) (models.Login, error) {

	query := `
		SELECT email, password, fullname, username  FROM users WHERE email = $1
	`
	var user models.Login

	rows, err := db.Conn.Query(context.Background(), query, email)
	if err != nil {
		return user, err
	}

	if !rows.Next() {
		return user, errors.New("user not found")
	}

	err = rows.Scan(&user.Email, &user.Password, &user.Fullname, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UserExist(email string) bool {

	query := `
		SELECT id FROM users 
		WHERE email = $1 LIMIT 1
	`

	var id int
	err := db.Conn.QueryRow(context.Background(), query, email).Scan(&id)

	return bool(err == nil)
}

func CreateUser(user models.DBUser) bool {

	query := `
		INSERT INTO users (fullname, username, email, password) 
		VALUES ($1, $2, $3, $4)
	`

	err := db.Conn.QueryRow(context.Background(), query, user.Fullname, user.Username, user.Email, user.Password).Scan()

	return err == nil
}
