package services

import (
	"github.com/CrowderSoup/socialboat/models"
	"github.com/jinzhu/gorm"
)

// PostService a service for getting, creating, and updating posts
type PostService struct {
	DB *gorm.DB
}

// NewPostService builds a new PostService
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{
		DB: db,
	}
}

// Create creates a new post in the database
func (s *PostService) Create(post *models.Post) error {
	if err := s.DB.Create(post).Error; err != nil {
		return err
	}

	return nil
}

// Update saves the post
func (s *PostService) Update(post *models.Post) error {
	if err := s.DB.Save(post).Error; err != nil {
		return err
	}

	return nil
}

// Delete delete's a post
func (s *PostService) Delete(post *models.Post) error {
	if err := s.DB.Delete(post).Error; err != nil {
		return err
	}

	return nil
}

// GetList gets a list of posts
func (s *PostService) GetList(page, limit int) ([]models.Post, error) {
	if limit == 0 {
		limit = 10
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	var posts []models.Post
	result := s.DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

// GetBySlug attempts to get a post by slug
func (s *PostService) GetBySlug(slug string) (*models.Post, error) {
	var post models.Post
	result := s.DB.Where(&models.Post{Slug: slug}).First(&post)
	if result.Error != nil {
		return nil, result.Error
	}

	return &post, nil
}
