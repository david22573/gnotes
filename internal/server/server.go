package server

import (
	"github.com/david22573/gnotes/internal/server/routes"
	"github.com/david22573/gnotes/internal/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = validators.NewCustomValidator()
	routes.RegisterAPIRoutes(e)

	return e
}
