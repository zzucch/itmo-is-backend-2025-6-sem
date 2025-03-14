package storage

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID    uint
	Message   string
	IsRead    bool
	CreatedAt time.Time
}
