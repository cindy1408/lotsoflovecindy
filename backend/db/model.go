package db

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ContentURL  string
	Description string
	Schedule    time.Time
	DateCreated time.Time
}
