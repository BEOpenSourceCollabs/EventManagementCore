package config_test

import (
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
)

func TestConfig_Environment(t *testing.T) {
	t.Run("check production env is production", func(t *testing.T) {
		if !config.Prod.IsProduction() {
			t.Error("expected IsProduction to return true")
		}
	})
	t.Run("check development env is not production", func(t *testing.T) {
		if config.Dev.IsProduction() {
			t.Error("expected IsProduction to return false")
		}
	})
	t.Run("check test env is not production", func(t *testing.T) {
		if config.Test.IsProduction() {
			t.Error("expected IsProduction to return false")
		}
	})
}
