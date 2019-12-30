package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/CrowderSoup/socialboat/config"
	"github.com/CrowderSoup/socialboat/controllers"
	"github.com/CrowderSoup/socialboat/migrations"
	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

func main() {
	// Load the config (pulls from config.toml first, then env variables for overrides)
	var s config.Server
	config.LoadConfig(&s, "config.toml")

	db, err := gorm.Open(s.DBConfig.Dialect, s.DBConfig.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(
		&models.Post{},
		&models.User{},
		&models.Profile{},
		&migrations.Migration{},
		&models.Menu{},
		&models.MenuItem{},
	)

	if s.Migrate {
		migrator, err := migrations.NewMigrator(db, s.MigrateUp)
		if err != nil {
			log.Fatal(err)
		}

		err = migrator.RunMigrations()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// Autoload Relationships
	db.Set("gorm:auto_preload", true)

	// Echo instance
	e := echo.New()
	e.Static(fmt.Sprintf("/%s", s.AssetsDir), s.AssetsDir)

	// Get our Session Store ready
	store := services.InitSessionStore(s.SessionSecret, db, true)
	e.Use(session.Middleware(store))

	// Custom Context
	ccHandler := controllers.NewCustomContextHandler(db, &s)
	e.Use(ccHandler.Handler)

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

	profileController := controllers.NewProfileController(db)
	profileController.InitRoutes(e.Group("/profile"))

	menuController := controllers.NewMenuController(db)
	menuController.InitRoutes(e.Group("/menus"))

	filesController := controllers.NewFilesController(db)
	filesController.InitRoutes(e.Group("/media"))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.Port)))
}
