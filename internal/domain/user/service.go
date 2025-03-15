package user

import "errors"

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateUser(user *User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return errors.New("username, email, and password are required")
	}
	return s.repository.Create(user)
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repository.FindByID(id)
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	return s.repository.FindByUsername(username)
}

func (s *Service) UpdateUser(user *User) error {
	if user.ID == 0 {
		return errors.New("invalid user ID")
	}
	return s.repository.Update(user)
}

func (s *Service) DeleteUser(id uint) error {
	return s.repository.Delete(id)
}
