package controllers

import (
	"database/sql"

	"github.com/leonibeldev/askme/db"
	"github.com/leonibeldev/askme/pkg/utils/models"
)

func GetUser(email string) (models.DBUser, error) {

	db := db.Connection()
	defer db.Close()

	query := `
		SELECT name, email, hashpassword, role, created_at FROM users WHERE email = ?
	`
	var user models.DBUser

	err := db.QueryRow(query, email).Scan(&user.Name, &user.Email, &user.HashPassword, &user.Role, &user.Created_at)
	if err == sql.ErrNoRows {
		return models.DBUser{}, sql.ErrNoRows
	}

	return user, nil
}

func UserExist(email string) bool {

	db := db.Connection()
	defer db.Close()

	query := `
		SELECT id FROM users WHERE email = ?
	`

	var id int
	err := db.QueryRow(query, email).Scan(&id)

	return bool(err == nil)
}

func CreateUser(user models.User) bool {

	db := db.Connection()
	defer db.Close()

	query := `
		INSERT INTO users (name, email, hashpassword) VALUES (?, ?, ?)
	`

	_, err := db.Exec(query, user.Name, user.Email, user.Password)

	if err == sql.ErrNoRows {
		return false
	}

	return err == nil
}
