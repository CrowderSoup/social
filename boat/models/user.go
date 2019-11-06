package models

import (
	"github.com/jinzhu/gorm"
)

// User the user of our application
type User struct {
	gorm.Model

	Nickname  string `gorm:"type:varchar(100);unique_index"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	Password  string
	FirstName string
	LastName  string
	PhotoURL  string
}
