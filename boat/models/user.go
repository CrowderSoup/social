package models

import (
	"github.com/jinzhu/gorm"
)

// User the user of our application
type User struct {
	gorm.Model

	Email    string `gorm:"type:varchar(100);unique_index"`
	Password string
}

// Profile the users profile
type Profile struct {
	gorm.Model

	UserID      int
	User        User
	NickName    string
	FirstName   string
	LastName    string
	PhotoURL    string
	PublicEmail string
	Phone       string
	Twitter     string
	Github      string
}
