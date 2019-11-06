package main

import (
	"github.com/foolin/goview/supports/echoview"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/CrowderSoup/social/boat/controllers"
	"github.com/CrowderSoup/social/boat/models"
)

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Post{})

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echoview.Default()

	postsController := controllers.NewPostsController(db)
	postsController.InitRoutes(e.Group("/"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
