package models

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DBUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	HashPassword string `json:"hashpassword"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
