package main

import (
	"github.com/foolin/goview"
	echoview "github.com/foolin/goview/supports/echoview-v4"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/CrowderSoup/social/boat/controllers"
	"github.com/CrowderSoup/social/boat/models"
	"github.com/CrowderSoup/social/boat/services"
)

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(
		&models.Post{},
		&models.User{},
		&models.Profile{},
	)

	// Echo instance
	e := echo.New()
	e.Static("/assets", "assets")

	// Get our Session Store ready
	store := services.InitSessionStore("secret", db, true)
	e.Use(session.Middleware(store))

	// Custom Context
	e.Use(controllers.CustomContextHandler)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echoview.New(goview.Config{
		Root:      "views",
		Extension: ".html",
		Master:    "layouts/master",
		Partials: []string{
			"partials/post",
			"partials/h-card",
		},
		DisableCache: true,
	})

	// HTTPErrorHandler
	e.HTTPErrorHandler = controllers.HTTPErrorHandler

	// Register Routes
	e.GET("/manifest.webmanifest", controllers.ManifestHandler)

	postsController := controllers.NewPostsController(db)
	postsController.InitRoutes(e.Group("/"))

	authController := controllers.NewAuthController(db)
	authController.InitRoutes(e.Group("/auth"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
