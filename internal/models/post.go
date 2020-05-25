package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Post a post
type Post struct {
	gorm.Model

	Title string
	Body  string
	Slug  string `gorm:"type:varchar(100);unique_index"`

	UserID int
	User   User
}

// FormattedDate returns the post's CreatedAt date, but formatted
func (p *Post) FormattedDate() string {
	return p.CreatedAt.Format(time.RFC822)
}
