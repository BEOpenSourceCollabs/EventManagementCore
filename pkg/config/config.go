package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
)

// Configuration .
type Configuration struct {
	Env      GoEnv
	Port     int
	Security SecurityConfiguration
	Database persist.DatabaseConfiguration
}

type SecurityConfiguration struct {
	JsonWebToken service.JsonWebTokenConfiguration
	Google       service.GoogleAuthenticationConfiguration
}

// NewEnvironmentConfiguration creates a configuration populated from os environment variables.
func NewEnvironmentConfiguration() Configuration {

	env, ok := os.LookupEnv("ENV")

	if !ok || env == "" {
		env = string(Dev)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		port = 8081
	}

	accessTokenSecret, ok := os.LookupEnv("ACCESS_TOKEN_SECRET")

	if !ok || accessTokenSecret == "" {
		panic(errors.New("invalid environment configuration: ACCESS_TOKEN_SECRET is required"))
	}

	refreshTokenSecret, ok := os.LookupEnv("REFRESH_TOKEN_SECRET")

	if !ok || refreshTokenSecret == "" {
		panic(errors.New("invalid environment configuration: REFRESH_TOKEN_SECRET is required"))
	}

	gClientId := os.Getenv("GOOGLE_CLIENT_ID")

	return Configuration{
		Port: port,
		Env:  ValidateEnv(GoEnv(env)),
		Security: SecurityConfiguration{
			JsonWebToken: service.JsonWebTokenConfiguration{
				AccessTokenSecret:  accessTokenSecret,
				RefreshTokenSecret: refreshTokenSecret,
			},
			Google: service.GoogleAuthenticationConfiguration{
				ClientId: gClientId,
			},
		},
		Database: persist.DatabaseConfiguration{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     os.Getenv("DATABASE_PORT"),
			Database: os.Getenv("DATABASE_NAME"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			SSLMode:  os.Getenv("DATABASE_SSL_MODE"),
		},
	}
}
