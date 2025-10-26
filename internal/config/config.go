package config

import (
	"echto/pkg/logger"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	APP_ENV  string `mapstructure:"APP_ENV"`
	APP_NAME string `mapstructure:"APP_NAME"`
	APP_PORT int    `mapstructure:"APP_PORT"`
	APP_HOST string `mapstructure:"APP_HOST"`
}

type DatabaseConfig struct {
	DB_HOST              string `mapstructure:"DB_HOST"`
	DB_PORT              int    `mapstructure:"DB_PORT"`
	DB_USER              string `mapstructure:"DB_USER"`
	DB_PASSWORD          string `mapstructure:"DB_PASSWORD"`
	DB_NAME              string `mapstructure:"DB_NAME"`
	DB_SSL_MODE          string `mapstructure:"DB_SSL_MODE"`
	DB_MAX_IDLE_CONNS    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_OPEN_CONNS    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_CONN_MAX_LIFETIME string `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

type LoggingConfig struct {
	LOG_LEVEL  string `mapstructure:"LOG_LEVEL"`
	LOG_FORMAT string `mapstructure:"LOG_FORMAT"`
}

type JWTConfig struct {
	JWT_SECRET       string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRE_HOURS int    `mapstructure:"JWT_EXPIRE_HOURS"`
}

func Load() *Config {
	// Set config file
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	// Set env prefix
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Log.Warn().Err(err).Msg("Failed to read config file, using defaults and environment variables")
	}

	var config Config = Config{
		App: AppConfig{
			APP_ENV:  viper.GetString("APP_ENV"),
			APP_NAME: viper.GetString("APP_NAME"),
			APP_PORT: viper.GetInt("APP_PORT"),
			APP_HOST: viper.GetString("APP_HOST"),
		},
		Database: DatabaseConfig{
			DB_HOST:              viper.GetString("DB_HOST"),
			DB_PORT:              viper.GetInt("DB_PORT"),
			DB_USER:              viper.GetString("DB_USER"),
			DB_PASSWORD:          viper.GetString("DB_PASSWORD"),
			DB_NAME:              viper.GetString("DB_NAME"),
			DB_SSL_MODE:          viper.GetString("DB_SSL_MODE"),
			DB_MAX_IDLE_CONNS:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			DB_MAX_OPEN_CONNS:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			DB_CONN_MAX_LIFETIME: viper.GetString("DB_CONN_MAX_LIFETIME"),
		},
		Logging: LoggingConfig{
			LOG_LEVEL:  viper.GetString("LOG_LEVEL"),
			LOG_FORMAT: viper.GetString("LOG_FORMAT"),
		},
		JWT: JWTConfig{
			JWT_SECRET:       viper.GetString("JWT_SECRET"),
			JWT_EXPIRE_HOURS: viper.GetInt("JWT_EXPIRE_HOURS"),
		},
	}

	return &config
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "echto")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "1h")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("JWT_EXPIRE_HOURS", 24)
	viper.SetDefault("APP_NAME", "echto")
	viper.SetDefault("APP_PORT", 9090)
	viper.SetDefault("APP_HOST", "localhost")
}
