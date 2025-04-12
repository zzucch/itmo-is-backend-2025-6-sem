package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

type Service struct {
	repository Repository
}

func (s *Service) AddToCart(userID uint, phoneID uint) error {
	_, err := s.repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	phone, err := s.repository.FindPhoneByID(phoneID)
	if err != nil || phone == nil {
		return errors.New("phone not found")
	}

	return s.repository.AddPhoneToCart(userID, phoneID)
}

func (s *Service) RemoveFromCart(userID uint, phoneID uint) error {
	_, err := s.repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.repository.RemovePhoneFromCart(userID, phoneID)
}

func (s *Service) GetCart(userID uint) ([]general.Phone, error) {
	user, err := s.repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user.Cart, nil
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

	updatedUser, err := s.repository.FindUserByID(user.ID)
	if err != nil {
		return "", err
	}
	*user = *updatedUser

	return tokenString, nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := s.repository.FindToken(tokenString)
	if err != nil {
		return nil, errors.New("token not found: " + err.Error())
	}
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("token expired")
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
		return nil, errors.New("validation: invalid token")
	}

	return claims, nil
}

func (s *Service) InvalidateToken(tokenString string) error {
	return s.repository.DeleteToken(tokenString)
}

func (s *Service) GetUserTokens(userID uint) ([]Token, error) {
	user, err := s.repository.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user.Tokens, nil
}

func (s *Service) InvalidateAllUserTokens(userID uint) error {
	user, err := s.repository.FindUserByID(userID)
	if err != nil {
		return err
	}

	for _, token := range user.Tokens {
		if err := s.repository.DeleteToken(token.Token); err != nil {
			return err
		}
	}

	return nil
}
