package dtos

import (
	"fmt"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/golang-jwt/jwt/v5"
)

type JwtPayload struct {
	Id   string
	Role constants.Role
}

// Sign generates a signed JWT token for given payload using the provided secret.
func (jwtPayload *JwtPayload) Sign(secret string) (string, error) {
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  jwtPayload.Id,
		"role": jwtPayload.Role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jwtPayload *JwtPayload) ParseSignedToken(tokenString, secret string) error {
	//parse JWT token using secret and same algorithm that is used to sign.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		(*jwtPayload) = JwtPayload{
			Id:   claims["sub"].(string),
			Role: constants.Role(claims["role"].(string)),
		}
	}

	return err
}
