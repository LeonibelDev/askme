package models

import (
	"time"
)

type DBUser struct {
	Id            int       `json:"id,omitempty"`
	Fullname      string    `json:"fullname"`
	Username      string    `json:"username,omitempty"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Role          string    `json:"role,omitempty"`
	Resume        string    `json:"resume,omitempty"`
	Is_verified   bool      `json:"is_verified,omitempty"`
	Created_at    time.Time `json:"created_at,omitempty"`
	Twitter       string    `json:"twitter"`
	Github        string    `json:"github"`
	Instagram     string    `json:"instagram"`
	External_link string    `json:"external_link"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GitHubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description"`
	Language    string `json:"language"`
}

type BlogPost struct {
	Position int    `json:"position" binding:"required"`
	Type     string `json:"Type" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type Post struct {
	ID       string     `json:"id,omitempty"`
	Title    string     `json:"title" binding:"required"`
	Cover    string     `json:"cover" binding:"required"`
	Author   string     `json:"author,omitempty"`
	Date     time.Time  `json:"date,omitempty"`
	Visible  bool       `json:"visible,omitempty"`
	Tags     []string   `json:"tags" binding:"required"`
	Sections []BlogPost `json:"sections,omitempty"`
}

type Newsletter struct {
	ID          string    `json:"id,omitempty"`
	Email       string    `json:"email"`
	Inserted_at time.Time `json:"Inserted_at,omitempty"`
}
