package db

import (
	"errors"
	"strings"

	"github.com/david22573/gnotes/internal/types/requests"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidPassword = errors.New("password must be at least 8 characters")
)

func (s *DBStore) CreateUser(req requests.CreateUserRequest) (*UserResponse, error) {
	// Validate input
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}
	if len(req.Password) < 8 {
		return nil, ErrInvalidPassword
	}

	// Check if email already exists
	exists, err := s.emailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: string(hashedPassword),
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *DBStore) GetUser(id uint) (*UserResponse, error) {
	var user User
	err := s.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *DBStore) GetUserByEmail(email string) (*UserResponse, error) {
	var user User
	err := s.DB.First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *DBStore) UpdateUser(id uint, updates map[string]interface{}) (*UserResponse, error) {
	var user User

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// First check if user exists
		if err := tx.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}

		// If email is being updated, validate it
		if email, ok := updates["email"].(string); ok {
			if err := validateEmail(email); err != nil {
				return err
			}
			// Check if new email already exists for different user
			exists, err := s.emailExistsExcludingID(email, id)
			if err != nil {
				return err
			}
			if exists {
				return ErrEmailExists
			}
			updates["email"] = strings.ToLower(strings.TrimSpace(email))
		}

		// If password is being updated, hash it
		if password, ok := updates["password"].(string); ok {
			if len(password) < 8 {
				return ErrInvalidPassword
			}
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			updates["password"] = string(hashedPassword)
		}

		// Perform the update
		return tx.Model(&user).Updates(updates).Error
	})

	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *DBStore) DeleteUser(id uint) error {
	var user User
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// First check if user exists
		if err := tx.First(&user, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}
		return tx.Delete(&user).Error
	})
	return err
}

func (s *DBStore) emailExists(email string) (bool, error) {
	var count int64
	err := s.DB.Model(&User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (s *DBStore) emailExistsExcludingID(email string, id uint) (bool, error) {
	var count int64
	err := s.DB.Model(&User{}).Where("email = ? AND id != ?", email, id).Count(&count).Error
	return count > 0, err
}

func validateEmail(email string) error {
	// Basic email validation - could be enhanced with regex or email validation package
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return ErrInvalidEmail
	}
	return nil
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
