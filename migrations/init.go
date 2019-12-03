package migrations

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/CrowderSoup/socialboat/models"
)

type initialMigration struct {
	name string
}

// NewInitialMigration returns the initialMigration
func NewInitialMigration() MigrationFile {
	return &initialMigration{
		name: "init",
	}
}

// Name returns the name
func (m *initialMigration) Name() string {
	return m.name
}

// Up runs the migration
func (m *initialMigration) Up(db *gorm.DB) error {
	// Post -> User
	if err := db.Debug().Model(&models.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		fmt.Println(err)
		return err
	}

	// Profile -> User
	if err := db.Model(&models.Profile{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}

	// MenuItem -> Menu
	if err := db.Model(&models.MenuItem{}).AddForeignKey("menu_id", "menus(id)", "RESTRICT", "RESTRICT").Error; err != nil {
		return err
	}

	return nil
}

// Down rolls back the migration
func (m *initialMigration) Down(db *gorm.DB) error {
	// Post -> User
	if err := db.Model(&models.Post{}).RemoveForeignKey("user_id", "users(id)").Error; err != nil {
		return err
	}

	// Profile -> User
	if err := db.Model(&models.Profile{}).RemoveForeignKey("user_id", "users(id)").Error; err != nil {
		return err
	}

	// MenuItem -> Menu
	if err := db.Model(&models.MenuItem{}).RemoveForeignKey("menu_id", "menus(id)").Error; err != nil {
		return err
	}

	return nil
}
