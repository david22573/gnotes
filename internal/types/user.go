package types

import (
	"errors"
	"fmt"
	"time"
)

// Constants and errors.
var ErrNameEmpty = errors.New("name cannot be empty")

// User struct (contains the password, but excluded from the JSON response).
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"-"` // Prevent password from being serialized.
	CreatedAt time.Time `json:"created_at"`
}

// Validate checks if the required fields are valid.
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("validation failed: %w", ErrNameEmpty)
	}
	return nil
}

// ToResponse converts a User to a UserResponse (excluding sensitive fields).
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}

// UserResponse struct defines the API response.
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest represents the JSON payload for user creation.
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
