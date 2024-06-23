package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hashes plaintext password using bcrypt algorithm
func HashPassword(plainText string) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashString := string(hashedPassword)
	return &hashString, nil
}

// compares a hashed password with its plaintext version using bcrypt algorithm.
func DoesPasswordMatch(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
