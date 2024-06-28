package service_test

import (
	"os"
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

var jwtService service.JsonWebTokenService

func init() {
	jwtService = service.NewJsonWebTokenService(
		&service.JsonWebTokenConfiguration{
			AccessTokenSecret:  "test123",
			RefreshTokenSecret: "test456",
		},
		logging.NewTextLogWriter(os.Stdout, logging.DEBUG),
	)
}

func TestJsonWebTokenService_InvalidInputs(t *testing.T) {
	t.Run("attempt sign access token with no payload", func(t *testing.T) {
		signed, err := jwtService.SignAccessToken(service.JwtPayload{})
		if err == nil {
			t.Error("jwt service should return error when empty jwt payload is provided")
		}
		if signed != nil {
			t.Error("jwt service should not return a signed string when empty jwt payload is provided")
		}
	})
	t.Run("attempt sign with invalid role", func(t *testing.T) {
		signed, err := jwtService.SignAccessToken(service.JwtPayload{Id: "test", Role: types.Role("sudo")})
		if err == nil {
			t.Error("jwt service should return error when empty jwt payload is provided")
		}
		if signed != nil {
			t.Error("jwt service should not return a signed string when empty jwt payload is provided")
		}
	})

	t.Run("attempt parse of invalid access token", func(t *testing.T) {
		parsed, err := jwtService.ParseAccessToken("")
		if err == nil {
			t.Error("jwt service should return error when attempting to parse empty string")
		}
		if parsed != nil {
			t.Error("jwt service should not return a jwt payload when attempting to parse empty string")
		}
	})

	t.Run("attempt sign refresh token with no payload", func(t *testing.T) {
		signed, err := jwtService.SignRefreshToken(service.RefreshTokenPayload{})
		if err == nil {
			t.Error("jwt service should return error when empty jwt payload is provided")
		}
		if signed != nil {
			t.Error("jwt service should not return a signed string when empty jwt payload is provided")
		}
	})

	t.Run("attempt parse of invalid refresh token", func(t *testing.T) {
		parsed, err := jwtService.ParseRefreshToken("")
		if err == nil {
			t.Error("jwt service should return error when attempting to parse empty string")
		}
		if parsed != nil {
			t.Error("jwt service should not return a jwt payload when attempting to parse empty string")
		}
	})

}

func TestJsonWebTokenService_SignAndParseSignedAccessToken(t *testing.T) {
	testPayload := service.JwtPayload{
		Id:   "test",
		Role: types.UserRole,
	}

	t.Run("sign and parse access token", func(t *testing.T) {
		signed, err := jwtService.SignAccessToken(testPayload)
		if err != nil {
			t.Error(err)
		}
		if signed == nil {
			t.Fatal("expected tokenString to be a valid jwt token string but was nil")
		}

		parsedPayload, err := jwtService.ParseAccessToken(*signed)
		if err != nil {
			t.Error(err)
		}
		if parsedPayload.Id != testPayload.Id {
			t.Errorf("expected id to be %s but was %s", testPayload.Id, parsedPayload.Id)
		}
		if parsedPayload.Role != testPayload.Role {
			t.Errorf("expected role to be %s but was %s", testPayload.Role, parsedPayload.Role)
		}
	})

}

func TestJsonWebTokenService_SignAndParseSignedRefreshToken(t *testing.T) {
	testPayload := service.RefreshTokenPayload{
		Id: "test",
	}

	t.Run("sign and parse refresh token", func(t *testing.T) {
		signed, err := jwtService.SignRefreshToken(testPayload)
		if err != nil {
			t.Error(err)
		}
		if signed == nil {
			t.Fatal("expected tokenString to be a valid jwt token string but was nil")
		}

		parsedPayload, err := jwtService.ParseRefreshToken(*signed)
		if err != nil {
			t.Error(err)
		}
		if parsedPayload.Id != testPayload.Id {
			t.Errorf("expected id to be %s but was %s", testPayload.Id, parsedPayload.Id)
		}
	})

}
