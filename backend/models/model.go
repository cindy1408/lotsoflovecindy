package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	ContentURL  string
	Description string
	UserID      uuid.UUID
	Schedule    time.Time
	DateCreated time.Time
}

type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name  string
	Email string
}
