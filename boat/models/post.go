package models

import "github.com/jinzhu/gorm"

// Post a post
type Post struct {
	gorm.Model

	Title string `json:"title"`
	Body  string `json:"body"`

	UserID int
	User   User
}
