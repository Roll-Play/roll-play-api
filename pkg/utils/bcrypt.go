package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(plaintext string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))

	if err != nil {
		return err
	}

	return nil
}
