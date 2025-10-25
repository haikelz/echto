package database

import (
	"echto/internal/config"
	echtoLogger "echto/pkg/logger"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg config.DatabaseConfig) *gorm.DB {
	// Build DSN
	dsn := buildDSN(cfg)

	// Configure GORM logger
	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Silent)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		echtoLogger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		echtoLogger.Log.Fatal().Err(err).Msg("Failed to get underlying sql.DB")
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// Parse connection max lifetime
	connMaxLifetime, err := time.ParseDuration(cfg.ConnMaxLifetime)
	if err != nil {
		echtoLogger.Log.Warn().Err(err).Str("value", cfg.ConnMaxLifetime).Msg("Invalid conn_max_lifetime, using default")
		connMaxLifetime = time.Hour
	}
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		echtoLogger.Log.Fatal().Err(err).Msg("Failed to ping database")
	}

	echtoLogger.Log.Info().Msg("Database connected successfully")
	return db
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
}
