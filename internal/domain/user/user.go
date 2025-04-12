package user

import (
	"errors"
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
	Tokens       []Token         `gorm:"foreignKey:UserID;references:ID"`
	Role
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func ValidateRole(role Role) error {
	switch role {
	case RoleUser, RoleAdmin:
		return nil
	default:
		return errors.New("invalid role")
	}
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

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
	Role     Role
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateResponse struct {
	Status string `json:"status"`
}

type AddToCartRequest struct {
	PhoneID uint `json:"phone_id" validate:"required"`
}

type RemoveFromCartRequest struct {
	PhoneID uint `json:"phone_id" validate:"required"`
}

type CartResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message,omitempty"`
	Phones  []general.Phone `json:"phones,omitempty"`
}
