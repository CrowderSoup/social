package models

// FileUploadReturn DTO for file upload
type FileUploadReturn struct {
	File File `json:"data"`
}

// File a file object
type File struct {
	FilePath string `json:"filePath"`
}
