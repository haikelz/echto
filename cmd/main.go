package main

import (
	"echto/internal/config"
	"echto/internal/database"
	"echto/internal/handler"
	"echto/internal/repository"
	routes "echto/internal/route"
	"echto/internal/service"
	"echto/pkg/logger"
	echtoMiddleware "echto/pkg/middleware"
	"fmt"

	_ "echto/pkg/swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// @title Echto API
// @version 1.0.0
// @description Clean Architecture API using Echo Framework
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@echto.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.Logging.LOG_LEVEL, cfg.Logging.LOG_FORMAT)

	// Initialize database
	db := database.Init(cfg.Database)

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Custom middleware
	e.Use(echtoMiddleware.RequestLogger())
	e.Use(middleware.RequestID())

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

	// Routes
	routes.UserRoute(e, userHandler)
	routes.SwaggerRoute(e)

	// Start server
	log.Info().Str("host", cfg.App.APP_HOST).Int("port", cfg.App.APP_PORT).Msg("Starting server")
	if err := e.Start(cfg.App.APP_HOST + ":" + fmt.Sprintf("%d", cfg.App.APP_PORT)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
