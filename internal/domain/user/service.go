package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateUser(user *User) error {
	if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("username, email, and password are required")
	}

	return s.repository.CreateUser(user)
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repository.FindUserByID(id)
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	return s.repository.FindUserByUsername(username)
}

func (s *Service) UpdateUser(user *User) error {
	if user.ID == 0 {
		return errors.New("invalid user ID")
	}

	return s.repository.UpdateUser(user)
}

func (s *Service) DeleteUser(id uint) error {
	return s.repository.DeleteUser(id)
}

func (s *Service) GetAllUsers() ([]*User, error) {
	return s.repository.GetAllUsers()
}

func (s *Service) Signup(req SignupRequest) (*User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.repository.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(req LoginRequest) (*User, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.repository.FindUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("no such user")
	}

	if !CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	user.LastLogin = time.Now()
	if err := s.repository.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GenerateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(tokenExpiry)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	userToken := Token{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}

	if err := s.repository.CreateToken(&userToken); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := s.repository.FindToken(tokenString)
	if err != nil || time.Now().After(token.ExpiresAt) {
		return nil, errors.New("invalid token")
	}

	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil || !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *Service) InvalidateToken(tokenString string) error {
	return s.repository.DeleteToken(tokenString)
}
