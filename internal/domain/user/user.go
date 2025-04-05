package user

import (
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash string
	LastLogin    time.Time
	Cart         []general.Phone `gorm:"many2many:user_cart;"`
	Tokens       []Token         `gorm:"foreignKey:UserID"`
}

type LoginRequest struct {
	Username string
	Password string
}

type SignupRequest struct {
	Username string
	Email    string
	Password string
}
