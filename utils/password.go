package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword ...
func HashPassword(password string) (string, error) {
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password, error: %w", err)
	}

	return string(hashedpassword), nil
}

// CheckPassword ...
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
