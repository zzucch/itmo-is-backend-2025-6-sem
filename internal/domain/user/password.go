package user

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	const bcryptCost = 12
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
