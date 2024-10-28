package routes

import (
	"github.com/david22573/gnotes/internal/handlers"
	"github.com/david22573/gnotes/internal/store/db"
	"github.com/labstack/echo/v4"
)

func RegisterAPIRoutes(e *echo.Echo, db *db.DBStore) {
	// Create handlers
	userHandler := handlers.NewUserHandler(db)

	// API routes
	api := e.Group("/api")

	// Users routes
	users := api.Group("/users")
	users.POST("", userHandler.CreateUser)
	users.GET("/:id", userHandler.GetUser)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}
