package server

import (
	"github.com/david22573/gnotes/internal/server/routes"
	"github.com/david22573/gnotes/internal/store/db"
	"github.com/david22573/gnotes/internal/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echo  *echo.Echo
	store *db.DBStore
}

func New(dbPath string) *Server {
	// Initialize the store
	db := db.NewDBStore(dbPath)

	// Create new echo instance
	e := echo.New()

	// Initialize server
	server := &Server{
		echo:  e,
		store: db,
	}

	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set up validator
	e.Validator = validators.NewCustomValidator()

	// Register routes
	routes.RegisterAPIRoutes(e, db)

	return server
}

// Start starts the server
func (s *Server) Start(address string) error {
	return s.echo.Start(address)
}
