package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hashes plaintext password using bcrypt algorithm
func HashPassword(plainText string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// compares a hashed password with its plaintext version using bcrypt algorithm.
func DoesPasswordMatch(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
