package models

import "github.com/jinzhu/gorm"

// Post a post
type Post struct {
	gorm.Model

	Title string
	Body  string
}
