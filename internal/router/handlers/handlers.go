package handlers

import (
	"github.com/labstack/echo/v4"
)

func GetNotes(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "Get all notes"})
}

func CreateNote(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "Note created"})
}
