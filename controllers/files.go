package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

// FilesController routes for uploading and managing files
type FilesController struct {
	FileService *services.FileService
}

// NewFilesController returns a new FilesController
func NewFilesController(db *gorm.DB) *FilesController {
	return &FilesController{
		FileService: services.NewFileService(db),
	}
}

// InitRoutes initialize routes for this controller
func (c *FilesController) InitRoutes(g *echo.Group) {
	g.POST("/upload", c.upload)
}

func (c *FilesController) upload(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	// Source
	file, err := bc.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Path
	uploadDir := fmt.Sprintf("assets/uploads/%s", time.Now().Format("2006/01/02"))
	uploadPath := fmt.Sprintf("%s/%s", uploadDir, file.Filename)

	// Create the uploadDir if it doesn't already exist
	err = createDirIfNotExist(uploadDir)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't create dir")
	}

	// Destination
	dst, err := os.Create(uploadPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't create file")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "couldn't save file")
	}

	f := models.File{
		FilePath: uploadPath,
	}

	location := fmt.Sprintf("%s/%s", bc.Server.RootURL, uploadPath)
	bc.Response().Header().Set("Location", location)

	return bc.JSON(http.StatusOK, models.FileUploadReturn{
		File: f,
	})
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
