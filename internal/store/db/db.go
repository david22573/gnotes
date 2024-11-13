package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBStore struct {
	DB *gorm.DB
}

func NewDBStore(dbPath string) *DBStore {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{}, &Note{})
	return &DBStore{
		DB: db,
	}
}
