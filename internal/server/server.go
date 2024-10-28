package server

import (
	"html/template"
	"io"
	"net/http"

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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var t = &Template{
	templates: template.Must(template.ParseGlob("public/views/**/*.html")),
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello.html", "World")
}

func New(dbPath string) *Server {
	// Initialize the store
	db := db.NewDBStore(dbPath)

	// Create new echo instance
	e := echo.New()

	e.Renderer = t
	e.GET("/hello", Hello)

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
