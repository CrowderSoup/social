package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"github.com/CrowderSoup/socialboat/models"
)

// UserService a service to handle dealing with Users
type UserService struct {
	DB *gorm.DB
}

// NewUserService returns a new UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

// Create creates a user using the given model
func (s *UserService) Create(user *models.User) error {
	var count int
	result := s.DB.Table("users").Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count >= 1 {
		return errors.New("only one user is allowed")
	}

	password, err := s.hashPassword(user.Password)
	if err != nil {
		return err
	}

	// Update model with the hashed version of the password before storage
	user.Password = password

	if err = s.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// GetByEmail gets a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword checks if a given password is correct for a user
func (s *UserService) CheckPassword(password string, user *models.User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}
