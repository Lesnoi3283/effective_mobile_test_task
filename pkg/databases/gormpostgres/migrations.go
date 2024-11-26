package gormpostgres

import (
	"fmt"
	"gorm.io/gorm"
	"musiclib/internal/app/entities"
)

// Migrate migrates entities from package "entities" to a database.
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(entities.Song{})
	if err != nil {
		return fmt.Errorf("failed to migrate migrations: %w", err)
	}
	return nil
}
