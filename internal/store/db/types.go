package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // This already includes ID uint, CreatedAt, UpdatedAt, and DeletedAt
	Name       string
	Email      string `gorm:"uniqueIndex"`
	Password   string
}
