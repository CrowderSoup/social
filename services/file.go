package services

import "github.com/jinzhu/gorm"

// FileService a service for managing files
type FileService struct {
	DB *gorm.DB
}

// NewFileService returns a new file service
func NewFileService(db *gorm.DB) *FileService {
	return &FileService{
		DB: db,
	}
}
