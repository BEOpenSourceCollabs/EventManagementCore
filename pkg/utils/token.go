package utils

import (
	"fmt"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/golang-jwt/jwt/v5"
)

// generates a signed JWT token for given payload using provided secret
func GenerateToken(payload *dtos.JwtPayload, secret string) (string, error) {
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  payload.Id,
		"role": payload.Role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// parses and validates a JWT token.
func ValidateToken(tokenString, secret string) (*dtos.JwtPayload, error) {
	//parse JWT token using secret and same algorithm that is used to sign.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload := &dtos.JwtPayload{
			Id:   claims["sub"].(string),
			Role: constants.Role(claims["role"].(string)),
		}
		return payload, nil
	}

	return nil, err
}
