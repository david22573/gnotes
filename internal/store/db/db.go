package db

import (
	"strconv"

	"github.com/david22573/gnotes/internal/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// UserModel represents the database model

// ToUser converts UserModel to types.User
func (m *User) ToUser() *types.User {
	return &types.User{
		ID:        strconv.FormatUint(uint64(m.ID), 10),
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	}
}

// FromUser converts types.User to UserModel
func (m *User) FromUser(user *types.User) {
	// Don't set ID when converting from User - let GORM handle it
	m.Name = user.Name
	m.Email = user.Email
	m.Password = user.Password
}

type DBStore struct {
	*gorm.DB
}

func NewDBStore(path string) *DBStore {
	if path == "" {
		path = "data.db"
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	return &DBStore{db}
}

// CreateUser creates a new user and returns the created user with ID
func (s *DBStore) CreateUser(req *types.CreateUserRequest) (*types.User, error) {
	model := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.DB.Create(model).Error
	if err != nil {
		return nil, err
	}

	return model.ToUser(), nil
}

func (s *DBStore) GetUser(id uint) (*types.User, error) {
	var model User
	err := s.DB.First(&model, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToUser(), nil
}

func (s *DBStore) UpdateUser(id uint, updates map[string]interface{}) (*types.User, error) {
	var model User
	err := s.DB.First(&model, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	err = s.DB.Model(&model).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return model.ToUser(), nil
}

func (s *DBStore) DeleteUser(id uint) error {
	result := s.DB.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return types.ErrUserNotFound
	}
	return nil
}

func (s *DBStore) GetUserByEmail(email string) (*types.User, error) {
	var model User
	err := s.DB.Where("email = ?", email).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToUser(), nil
}
