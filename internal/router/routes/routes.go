package routes

import (
	"github.com/david22573/gnotes/internal/router/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterAPIRoutes(e *echo.Echo) {
	api := e.Group("/api")

	v1 := api.Group("/v1")
	{
		notes := v1.Group("/notes")
		notes.GET("", handlers.GetNotes)
		notes.POST("", handlers.CreateNote)
	}
}
