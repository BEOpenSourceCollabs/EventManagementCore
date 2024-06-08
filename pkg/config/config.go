package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
)

// Configuration .
type Configuration struct {
	Secret   string
	Env      GoEnv
	Port     int
	Database persist.DatabaseConfiguration
}

// NewEnvironmentConfiguration creates a configuration populated from os environment variables.
func NewEnvironmentConfiguration() Configuration {

	env := os.Getenv("ENV")

	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		port = 8081
	}

	secret, ok := os.LookupEnv("SECRET")

	if !ok || secret == "" {
		panic(errors.New("invalid environment configuration: SECRET is required"))
	}

	return Configuration{
		Port:   port,
		Env:    ValidateEnv(GoEnv(env)),
		Secret: secret,
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
