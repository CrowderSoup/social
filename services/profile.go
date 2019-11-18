package services

import (
	"github.com/CrowderSoup/socialboat/models"
	"github.com/jinzhu/gorm"
)

// ProfileService a service to manage profiles
type ProfileService struct {
	DB *gorm.DB
}

// NewProfileService returns a new ProfileService
func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{
		DB: db,
	}
}

// Create saves a new profile to the database
func (s *ProfileService) Create(profile *models.Profile) error {
	if err := s.DB.Create(profile).Error; err != nil {
		return err
	}

	return nil
}

// Update updates a profile
func (s *ProfileService) Update(profile *models.Profile) error {
	if err := s.DB.Save(&profile).Error; err != nil {
		return err
	}

	return nil
}

// GetFirst gets the first profile from the db (there should only be one)
func (s *ProfileService) GetFirst() (*models.Profile, error) {
	var profile models.Profile

	if err := s.DB.First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

// GetByUserID gets a profile by user
func (s *ProfileService) GetByUserID(userID int) (*models.Profile, error) {
	var profile models.Profile

	// Get Profile and check for error
	if err := s.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}
