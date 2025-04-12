package user

import "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"

type Repository interface {
	CreateUser(user *User) error
	FindUserByID(id uint) (*User, error)
	FindUserByUsername(username string) (*User, error)
	GetAllUsers() ([]*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	CreateToken(token *Token) error
	FindToken(tokenString string) (*Token, error)
	DeleteToken(tokenString string) error
	DeleteExpiredTokens() error
	RemovePhoneFromCart(userID uint, phoneID uint) error
	AddPhoneToCart(userID uint, phoneID uint) error
	FindPhoneByID(id uint) (*general.Phone, error)
}
