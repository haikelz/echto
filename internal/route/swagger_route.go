package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SwaggerRoute(e *echo.Echo) {
	echoSwagger.InstanceName("Echto API Documentation")
	echoSwagger.SyntaxHighlight(true)

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler(
		echoSwagger.PersistAuthorization(true),
	))
}
