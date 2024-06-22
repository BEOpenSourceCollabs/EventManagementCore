package service

import (
	"fmt"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
	"github.com/golang-jwt/jwt/v5"
)

// JwtPayload represents the json web token claims to be signed or that have been parsed.
type JwtPayload struct {
	Id   string
	Role types.Role
}

// JsonWebTokenService for signing and parsing json web tokens.
type JsonWebTokenService interface {
	Sign(JwtPayload) (*string, error)
	ParseSignedToken(string) (*JwtPayload, error)
}

// JsonWebTokenConfiguration settings for the json web token service.
type JsonWebTokenConfiguration struct {
	Secret string
}

type jsonWebTokenService struct {
	logger logging.Logger
	config *JsonWebTokenConfiguration
}

// NewJsonWebTokenService creates a new implementation of the JsonWebTokenService.
func NewJsonWebTokenService(config *JsonWebTokenConfiguration, lw logging.LogWriter) JsonWebTokenService {
	return &jsonWebTokenService{
		logger: logging.NewContextLogger(lw, "JsonWebTokenService"),
		config: config,
	}
}

// Sign generates a signed JWT token for given payload using the provided secret.
func (svc *jsonWebTokenService) Sign(payload JwtPayload) (*string, error) {
	// Create a new token object, specifying signing method and claims

	svc.logger.Debugf("signing payload with sub: '%s', role: '%s'", payload.Id, payload.Role)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  payload.Id,
		"role": payload.Role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(svc.config.Secret))
	svc.logger.Debugf("signing payload: '%s'", tokenString)

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func (svc *jsonWebTokenService) ParseSignedToken(tokenString string) (*JwtPayload, error) {
	//parse JWT token using secret and same algorithm that is used to sign.

	svc.logger.Debugf("parsing signed token: '%s'", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(svc.config.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		svc.logger.Debugf("parsing token claims: id => '%s', role => '%s'", claims["sub"], claims["role"])
		return &JwtPayload{
			Id:   claims["sub"].(string),
			Role: types.Role(claims["role"].(string)),
		}, nil
	}

	return nil, err
}
