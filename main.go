package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/koding/multiconfig"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/CrowderSoup/socialboat/controllers"
	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

// Server our server
type Server struct {
	AssetsDir      string `default:"assets"`
	DBConfig       DBConfig
	Port           int `default:"8080"`
	RendererConfig services.RendererConfig
	SessionSecret  string `required:"true"`
}

// DBConfig Config for the database
type DBConfig struct {
	ConnectionString string `default:"boat.db"`
	Dialect          string `default:"sqlite3"`
}

func main() {
	// Load the config (pulls from config.toml first, then env variables for overrides)
	var s Server
	m := multiconfig.NewWithPath("config.toml")
	m.MustLoad(&s)

	db, err := gorm.Open(s.DBConfig.Dialect, s.DBConfig.ConnectionString)
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
	e.Static(fmt.Sprintf("/%s", s.AssetsDir), s.AssetsDir)

	// Get our Session Store ready
	store := services.InitSessionStore(s.SessionSecret, db, true)
	e.Use(session.Middleware(store))

	// Custom Context
	e.Use(controllers.CustomContextHandler)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = services.NewRenderer(s.RendererConfig)

	// HTTPErrorHandler
	e.HTTPErrorHandler = controllers.HTTPErrorHandler

	// Register Routes
	e.GET("/manifest.webmanifest", controllers.ManifestHandler)

	postsController := controllers.NewPostsController(db)
	postsController.InitRoutes(e.Group("/"))

	authController := controllers.NewAuthController(db)
	authController.InitRoutes(e.Group("/auth"))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.Port)))
}
