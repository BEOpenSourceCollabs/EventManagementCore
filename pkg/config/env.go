package config

import "strings"

type GoEnv string

func (env GoEnv) IsProduction() bool {
	return strings.EqualFold(string(env), "production")
}

const (
	Dev  GoEnv = "development"
	Prod GoEnv = "production"
	Test GoEnv = "test"
)

func ValidateEnv(env GoEnv) GoEnv {
	switch env {
	case Dev, Prod, Test:
		return env
	default:
		return Dev
	}
}
