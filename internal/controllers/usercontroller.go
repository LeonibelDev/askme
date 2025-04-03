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
		SELECT name, email, hashpassword, role, created_at FROM users WHERE email = $1
	`
	var user models.DBUser

	rows, err := db.Conn.Query(context.Background(), query, email)
	if err != nil {
		return user, err
	}

	if !rows.Next() {
		return user, errors.New("user not found")
	}

	err = rows.Scan(&user.Name, &user.Email, &user.HashPassword, &user.Role, &user.Created_at)
	if err != nil {
		return user, err
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

func CreateUser(user models.User) bool {

	db.DataBaseConn()

	query := `
		INSERT INTO users (name, email, hashpassword) 
		VALUES ($1, $2, $3)
	`

	err := db.Conn.QueryRow(context.Background(), query, user.Name, user.Email, user.Password).Scan()

	return err == nil || err == pgx.ErrNoRows
}