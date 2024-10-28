package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/david22573/gnotes/internal/store/db"
	"github.com/david22573/gnotes/internal/types"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store *db.DBStore
}

func NewUserHandler(store *db.DBStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// @Summary Create a new user
// @Description Creates a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.CreateUserRequest true "User creation request"
// @Success 201 {object} types.UserResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req types.CreateUserRequest

	// Bind request body to struct
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Validate the request struct
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Check if user already exists by email
	existingUser, err := h.store.GetUserByEmail(req.Email)
	if err != nil && !errors.Is(err, types.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking user existence")
	}
	if existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "Email already registered")
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing password")
	}

	// Create user
	user, err := h.store.CreateUser(&types.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	return c.JSON(http.StatusCreated, user.ToResponse())
}

// @Summary Get a user by ID
// @Description Retrieves a user's information by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} types.UserResponse
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.store.GetUser(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error finding user")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

// @Summary Update a user
// @Description Updates a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body types.UpdateUserRequest true "User update request"
// @Success 200 {object} types.UserResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// You'll need to create this type
	var req types.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error processing password")
		}
		updates["password"] = hashedPassword
	}

	user, err := h.store.UpdateUser(uint(id), updates)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating user")
	}

	return c.JSON(http.StatusOK, user.ToResponse())
}

// @Summary Delete a user
// @Description Deletes a user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.store.DeleteUser(uint(id))
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting user")
	}

	return c.NoContent(http.StatusNoContent)
}

// Helper functions
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedBytes), nil
}

// Add this to your types package
type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=50,username"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,containsany=!@#$%^&*"`
}
