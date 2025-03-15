package user

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string // hehe
	IsAdmin  bool
	Cart     []general.Phone
}
