package user

type Repository interface {
	CreateUser(user *User) error
	FindUserByID(id uint) (*User, error)
	FindUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
}
