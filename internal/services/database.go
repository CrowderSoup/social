package services

import (
	"context"

	"github.com/CrowderSoup/socialboat/internal/config"
	"go.uber.org/fx"

	"github.com/jinzhu/gorm"
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
	return &Database{
		config: c.DBConfig,
	}
}

// InvokeDB opens / manages our database connection
func InvokeDB(lc fx.Lifecycle, d *Database) {
	lc.Append(
		fx.Hook{
			OnStart: func(context context.Context) error {
				db, err := gorm.Open(d.config.Dialect, d.config.ConnectionString)
				if err != nil {
					return err
				}

				// Save our connection for later use
				d.connection = db
				return nil
			},
			OnStop: func(context context.Context) error {
				d.connection.Close()
				return nil
			},
		},
	)
}
