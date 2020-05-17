package services

import (
	"context"

	"github.com/CrowderSoup/socialboat/internal/config"
	"github.com/CrowderSoup/socialboat/migrations"
	"github.com/CrowderSoup/socialboat/models"

	"github.com/jinzhu/gorm"
	"go.uber.org/fx"

	// Driver for gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database holds our gorm db
type Database struct {
	connection *gorm.DB
	config     config.DBConfig
}

// NewDatabase returns a Database
func NewDatabase(c *config.Config) *Database {
	db, err := gorm.Open(c.DBConfig.Dialect, c.DBConfig.ConnectionString)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.Post{},
		&models.User{},
		&models.Profile{},
		&models.Menu{},
		&models.MenuItem{},
		&migrations.Migration{},
	)

	// Autoload Relationships
	db.Set("gorm:auto_preload", true)

	return &Database{
		config:     c.DBConfig,
		connection: db,
	}
}

// InvokeDB opens / manages our database connection
func InvokeDB(lc fx.Lifecycle, d *Database) {
	lc.Append(
		fx.Hook{
			OnStart: func(context context.Context) error {
				return nil
			},
			OnStop: func(context context.Context) error {
				// Ensure we close our connection on app shutdown
				d.connection.Close()
				return nil
			},
		},
	)
}
