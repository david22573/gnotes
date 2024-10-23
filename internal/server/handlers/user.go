package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/david22573/gnotes/internal/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// UserStore defines the interface for user storage operations
type UserStore interface {
	Create(user types.User) error
	FindByName(name string) (*types.User, error)
}

// UserHandler handles HTTP requests related to users
type UserHandler struct {
	store UserStore
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(store UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// CreateUser handles the creation of new users
// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param request body types.CreateUserRequest true "User creation request"
// @Success 201 {object} types.UserResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Router /users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req types.CreateUserRequest

	// Bind request body to struct
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	// Validate the request struct
	if err := c.Validate(&req); err != nil {
		return err // Our custom validator already returns proper HTTP errors
	}

	// Check if user already exists
	existingUser, err := h.store.FindByName(req.Name)
	if err != nil && !errors.Is(err, types.ErrUserNotFound) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error checking user existence")
	}
	if existingUser != nil {
		return echo.NewHTTPError(http.StatusConflict, "Username already taken")
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error processing password")
	}

	// Create user
	newUser := types.User{
		ID:        generateUUID(),
		Name:      req.Name,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.store.Create(newUser); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	return c.JSON(http.StatusCreated, newUser.ToResponse())
}

// validateCreateUserRequest validates the user creation request
func validateCreateUserRequest(req types.CreateUserRequest) error {
	if len(req.Name) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	// Add more validation rules as needed
	return nil
}

// hashPassword securely hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedBytes), nil
}

// generateUUID generates a unique identifier
func generateUUID() string {
	return uuid.New().String()
}

// InMemoryUserStore implements UserStore interface for development/testing
type InMemoryUserStore struct {
	users []types.User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make([]types.User, 0),
	}
}

func (s *InMemoryUserStore) Create(user types.User) error {
	s.users = append(s.users, user)
	return nil
}

func (s *InMemoryUserStore) FindByName(name string) (*types.User, error) {
	for _, user := range s.users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, types.ErrUserNotFound
}
