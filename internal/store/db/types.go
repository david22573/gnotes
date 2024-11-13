package db

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     *User `gorm:"foreignKey:UserID"`
	IsPublic bool
	Created  string
}

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Notes    []*Note `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID    uint
	Name  string
	Email string
}
