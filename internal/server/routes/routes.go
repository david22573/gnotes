package routes

import (
	"github.com/david22573/gnotes/internal/server/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")
	v1 := api.Group("/v1")
	registerUserRoutes(v1)
}

func registerUserRoutes(v1 *echo.Group) {
	v1.POST("/users", handlers.CreateUser)
}
