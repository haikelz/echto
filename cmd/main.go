package main

import (
	"echto/internal/config"
	"echto/internal/database"
	"echto/internal/handler"
	"echto/internal/repository"
	"echto/internal/service"
	"echto/pkg/logger"
	echtoMiddleware "echto/pkg/middleware"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)

	// Initialize database
	db := database.Init(cfg.Database)

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService)

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

	// Routes
	api := e.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"service": "echto",
		})
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	log.Info().Str("host", cfg.App.Host).Int("port", cfg.App.Port).Msg("Starting server")
	if err := e.Start(cfg.App.Host + ":" + fmt.Sprintf("%d", cfg.App.Port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
