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

// RefreshTokenPayload represents json web token claims to be signed or that have been parsed.
type RefreshTokenPayload struct {
	Id string
}

// JsonWebTokenService for signing and parsing json web tokens.
type JsonWebTokenService interface {
	SignAccessToken(JwtPayload) (*string, error)
	ParseAccessToken(string) (*JwtPayload, error)
	SignRefreshToken(RefreshTokenPayload) (*string, error)
	ParseRefreshToken(string) (*RefreshTokenPayload, error)
}

// JsonWebTokenConfiguration settings for the json web token service.
type JsonWebTokenConfiguration struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
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

// sign generates a signed JWT token for given payload using the provided secret.
func sign(claims jwt.MapClaims, secret string) (*string, error) {
	// Create a new token object, specifying signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

// parse validates signature of JWT token against provided secret, parses payload if successful and returns claims
func parse(tokenString string, secret string) (jwt.MapClaims, error) {
	//parse JWT token using secret and same algorithm that is used to sign.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (svc *jsonWebTokenService) SignAccessToken(payload JwtPayload) (*string, error) {

	if len(payload.Id) < 1 {
		return nil, fmt.Errorf("jwtPayload must contain an id")
	}
	if !payload.Role.IsValid() {
		return nil, fmt.Errorf("jwtPayload must contain a valid role")
	}

	claims := jwt.MapClaims{
		"sub":  payload.Id,
		"role": payload.Role,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	svc.logger.Debugf("signing acccess token payload with sub: '%s', role: '%s'", payload.Id, payload.Role)
	token, err := sign(claims, svc.config.AccessTokenSecret)
	svc.logger.Debugf("signed access token payload: '%s'", *token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (svc *jsonWebTokenService) ParseAccessToken(tokenString string) (*JwtPayload, error) {

	svc.logger.Debugf("parsing signed access token: '%s'", tokenString)
	claims, err := parse(tokenString, svc.config.AccessTokenSecret)

	if err != nil {
		return nil, err
	}

	svc.logger.Debugf("parsed access token claims: id => '%s', role => '%s'", claims["sub"], claims["role"])

	return &JwtPayload{
		Id:   claims["sub"].(string),
		Role: types.Role(claims["role"].(string)),
	}, nil
}

func (svc *jsonWebTokenService) SignRefreshToken(payload RefreshTokenPayload) (*string, error) {
	if len(payload.Id) < 1 {
		return nil, fmt.Errorf("RefreshTokenPayload must contain an id")
	}

	week := time.Hour * 24 * 7

	claims := jwt.MapClaims{
		"sub": payload.Id,
		"exp": time.Now().Add(week).Unix(), //7 day expiry
	}

	svc.logger.Debugf("signing refresh token payload with sub: '%s'", payload.Id)
	token, err := sign(claims, svc.config.RefreshTokenSecret)
	svc.logger.Debugf("signed refresh token payload: '%s'", *token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (svc *jsonWebTokenService) ParseRefreshToken(token string) (*RefreshTokenPayload, error) {
	svc.logger.Debug("validating access token")

	claims, err := parse(token, svc.config.RefreshTokenSecret)

	if err != nil {
		return nil, err
	}

	svc.logger.Debugf("parsing refresh token claims: id => '%s'", claims["sub"])

	return &RefreshTokenPayload{
		Id: claims["sub"].(string),
	}, nil
}
