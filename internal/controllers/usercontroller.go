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
		SELECT fullname, email, password, role, created_at FROM users WHERE email = $1
	`
	var user models.DBUser

	rows, err := db.Conn.Query(context.Background(), query, email)
	if err != nil {
		return user, err
	}

	if !rows.Next() {
		return user, errors.New("user not found")
	}

	err = rows.Scan(&user.Fullname, &user.Email, &user.Password, &user.Role, &user.Created_at)
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

func CreateUser(user models.DBUser) bool {

	db.DataBaseConn()

	query := `
		INSERT INTO users (fullname, username, email, password) 
		VALUES ($1, $2, $3, $4)
	`

	err := db.Conn.QueryRow(context.Background(), query, user.Fullname, user.Username, user.Email, user.Password).Scan()

	return err == nil || err == pgx.ErrNoRows
}
