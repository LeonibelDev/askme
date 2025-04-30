package controllers

import (
	"context"

	"github.com/leonibeldev/askme/db"
)

func AddUserNewsletter(email string) error {
	db.DataBaseConn()

	query := `
		 INSERT INTO newsletter(email)
		 VALUES ($1)
	`

	_, err := db.Conn.Exec(context.Background(), query, email)
	if err != nil {
		return err
	}

	defer db.Conn.Close(context.Background())

	return nil
}

func RemoveUserNewsletter(uuid string) error {

	db.DataBaseConn()

	query := `
		DELETE FROM newsletter
		WHERE id = $1
	`

	_, err := db.Conn.Exec(context.Background(), query, uuid)
	if err != nil {
		return err
	}

	return nil

}
