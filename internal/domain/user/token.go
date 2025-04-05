package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

const (
	tokenExpiry = 24 * time.Hour
	secretKey   = "very-secret-key-hehe"
)

type Token struct {
	gorm.Model
	UserID    uint
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}
