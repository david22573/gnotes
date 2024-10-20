package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/david22573/gnotes/internal/types"
	"github.com/labstack/echo/v4"
)

// In-memory user store for simplicity (replace with DB in production).
var users = []types.User{}

// Create a new user.
func CreateUser(c echo.Context) error {
	var req types.CreateUserRequest

	// Bind and validate the request body.
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation failed: "+err.Error())
	}

	// Create the user object with hashed password.
	newUser := types.User{
		ID:        generateUserID(),
		Name:      req.Name,
		Password:  hashPassword(req.Password), // Hashing the password for security.
		CreatedAt: time.Now(),
	}

	// Store the user in memory.
	users = append(users, newUser)

	// Respond with the safe user response.
	return c.JSON(http.StatusCreated, newUser.ToResponse())
}

// Generate a unique user ID (for demo purposes).
func generateUserID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Hash the password (mocked for demo; use bcrypt in production).
func hashPassword(password string) string {
	return "hashed_" + password // Replace with bcrypt or another secure hash.
}
