package main

import (
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

	// Get our Session Store ready
	store := services.InitSessionStore("secret", db, true)
	e.Use(session.Middleware(store))

	// Custom Context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s, err := services.GetSession("Boat", c)
			if err != nil {
				return err
			}

			cc := &controllers.BoatContext{
				Context: c,
				Session: s,
			}
			return next(cc)
		}
	})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echoview.Default()

	// HTTPErrorHandler
	e.HTTPErrorHandler = controllers.HTTPErrorHandler

	// Register Routes
	postsController := controllers.NewPostsController(db)
	postsController.InitRoutes(e.Group("/"))

	authController := controllers.NewAuthController(db)
	authController.InitRoutes(e.Group("/auth"))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
