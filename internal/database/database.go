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
	sqlDB.SetMaxIdleConns(cfg.DB_MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(cfg.DB_MAX_OPEN_CONNS)

	// Parse connection max lifetime
	connMaxLifetime, err := time.ParseDuration(cfg.DB_CONN_MAX_LIFETIME)
	if err != nil {
		echtoLogger.Log.Warn().Err(err).Str("value", cfg.DB_CONN_MAX_LIFETIME).Msg("Invalid conn_max_lifetime, using default")
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
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_NAME,
		cfg.DB_SSL_MODE,
	)
}
