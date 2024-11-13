package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/david22573/gnotes/internal/store/db"
	"github.com/david22573/gnotes/internal/types/requests"
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

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req requests.CreateUserRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing password")
	}

	user, err := h.store.CreateUser(requests.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.store.GetUser(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var req requests.UpdateUserRequest
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
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	err = h.store.DeleteUser(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting user")
	}

	return c.NoContent(http.StatusNoContent)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedBytes), nil
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=3,max=50,username"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8,containsany=!@#$%^&*"`
}
