package models

import "github.com/jinzhu/gorm"

// Post a post
type Post struct {
	gorm.Model

	Title string
	Body  string
	Slug  string `gorm:"type:varchar(100);unique_index"`

	UserID int
	User   User
}
