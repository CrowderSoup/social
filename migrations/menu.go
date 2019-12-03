package migrations

import (
	"github.com/CrowderSoup/socialboat/models"
	"github.com/jinzhu/gorm"
)

type menuMigration struct {
	name string
}

// NewMenuMigration returns a new menu migration
func NewMenuMigration() MigrationFile {
	return &menuMigration{
		name: "menu",
	}
}

func (m *menuMigration) Name() string {
	return m.name
}

// Up runs the migration
func (m *menuMigration) Up(db *gorm.DB) error {
	menu := models.Menu{
		Name: "Default",
	}

	if err := db.Create(&menu).Error; err != nil {
		return err
	}

	return nil
}

// Down rolls back the migration
func (m *menuMigration) Down(db *gorm.DB) error {
	// Do nothing, we don't want to destroy data
	return nil
}
