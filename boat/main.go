package main

import (
	"github.com/foolin/goview/supports/echoview"
	session "github.com/ipfans/echo-session"
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

	db.AutoMigrate(
		&models.Post{},
		&models.User{},
		&models.Profile{},
	)

	// Echo instance
	e := echo.New()
	store := session.NewCookieStore([]byte("secret"))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Sessions("GSESSION", store))

	//Set Renderer
	e.Renderer = echoview.Default()

	// Register Routes
	postsController := controllers.NewPostsController(db)
	postsController.InitRoutes(e.Group("/"))

	authController := controllers.NewAuthController(db)
	authController.InitRoutes(e.Group("/auth"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
