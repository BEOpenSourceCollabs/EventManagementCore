package config

type GoEnv string

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
