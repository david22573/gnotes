package types

import (
	"errors"
	"time"
)

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (n Note) Validate() error {
	if n.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}
