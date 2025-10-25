package database

import (
	"echto/internal/entity"
	"echto/pkg/logger"

	"gorm.io/gorm"
)

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	logger.Log.Info().Msg("Running database migrations...")

	// Auto migrate all entities
	if err := db.AutoMigrate(
		&entity.User{},
	); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to run auto migration")
		return err
	}

	logger.Log.Info().Msg("Database migrations completed successfully")
	return nil
}
