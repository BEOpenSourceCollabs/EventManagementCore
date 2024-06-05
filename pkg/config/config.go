package config

import (
	"os"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
)

// Configuration .
type Configuration struct {
	Database persist.DatabaseConfiguration
}

// NewEnvironmentConfiguration creates a configuration populated from os environment variables.
func NewEnvironmentConfiguration() Configuration {
	return Configuration{
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
