package middleware

import (
	"echto/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logger.Log.Info().
				Str("method", c.Request().Method).
				Str("uri", values.URI).
				Int("status", values.Status).
				Dur("latency", values.Latency).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", c.Request().UserAgent()).
				Msg("HTTP request")
			return nil
		},
	})
}

// RequestID returns a middleware that generates a unique request ID
func RequestID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return generateRequestID()
		},
	})
}

// RateLimiter returns a rate limiter middleware
func RateLimiter(store middleware.RateLimiterStore) echo.MiddlewareFunc {
	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
		DenyHandler: func(c echo.Context, identifier string, err error) error {
			return c.JSON(429, map[string]string{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests",
			})
		},
	})
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
