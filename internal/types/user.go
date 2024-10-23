package types

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

// Custom validator instance
var (
	validate        = validator.New()
	ErrUserNotFound = errors.New("user not found")
)

// CreateUserRequest represents the request body for user creation
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50,username"`
	Password string `json:"password" validate:"required,min=8,containsany=!@#$%^&*"`
	Email    string `json:"email" validate:"required,email"`
}

// UserResponse represents the safe user data to return
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// User represents the internal user model
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never send password in response
	CreatedAt time.Time `json:"created_at"`
}

// Custom validation rules
func init() {
	// Register custom username validator
	validate.RegisterValidation("username", validateUsername)
}

// validateUsername ensures username contains only allowed characters
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	// Allow alphanumeric characters and underscores
	for _, char := range username {
		if !isAlphanumeric(char) && char != '_' {
			return false
		}
	}
	return true
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9')
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
