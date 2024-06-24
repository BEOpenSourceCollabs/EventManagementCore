package utils_test

import (
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

func TestUtils_PasswordHashAndCheck(t *testing.T) {
	pw := "T3sTpa$$w0rd123!!45"

	hash, err := utils.HashPassword(pw)
	if err != nil {
		t.Fatalf("expected no error while hashing password but got %v", err)
	}
	if hash == nil {
		t.Fatal("expected hash password to not be nil")
	}
	if !utils.DoesPasswordMatch(pw, *hash) {
		t.Errorf("expected DoesPasswordsMatch to return true when given hash '%s' and the original password '%s'", *hash, pw)
	}
}
