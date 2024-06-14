package routes_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/routes"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/test/mock"
)

var userContextHelper routes.UserContextHelpers

func init() {
	// Create a mock user repository, mocking only the required GetUserByID function.
	// Hardcode a user for each of the roles, for later texting of the LoadUserFromContextWithRole member function of the UserContextHelper
	var mockUserRepository repository.UserRepository = mock.UserRepository{
		GetUserByIDFn: func(id string) (*models.UserModel, error) {
			switch id {
			case "user":
				return &models.UserModel{
					ID:        id,
					Username:  "test1",
					Email:     "test1@domain.com",
					FirstName: sql.NullString{String: "unit1", Valid: true},
					LastName:  sql.NullString{String: "test1", Valid: true},
					Role:      constants.UserRole,
				}, nil
			case "admin":
				return &models.UserModel{
					ID:        id,
					Username:  "test2",
					Email:     "test2@domain.com",
					FirstName: sql.NullString{String: "unit2", Valid: true},
					LastName:  sql.NullString{String: "test2", Valid: true},
					Role:      constants.AdminRole,
				}, nil
			case "organizer":
				return &models.UserModel{
					ID:        id,
					Username:  "test3",
					Email:     "test3@domain.com",
					FirstName: sql.NullString{String: "unit3", Valid: true},
					LastName:  sql.NullString{String: "test3", Valid: true},
					Role:      constants.OrganizerRole,
				}, nil
			}
			return nil, fmt.Errorf("no user found with id %s", id)
		},
	}

	// Initialize the user context helper with the mock user repository, that will be tested.
	userContextHelper = routes.UserContextHelpers{
		&mockUserRepository,
	}
}

func TestUserContextHelpers_LoadUserFromContext(t *testing.T) {
	t.Run("User exists", func(t *testing.T) {
		payload := &dtos.JwtPayload{
			Id:   "user",
			Role: "user",
		}
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContext(r.WithContext(
			context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, payload),
		))
		if err != nil {
			t.Errorf(err.Error())
		}
		if user == nil {
			t.Errorf("failed to load mock user")
		}
	})

	t.Run("User doesn't exist", func(t *testing.T) {
		payload := &dtos.JwtPayload{
			Id:   "test",
			Role: "user",
		}
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContext(r.WithContext(
			context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, payload),
		))
		if err == nil {
			t.Errorf("expected error when attempting to load non-existent user from context")
		}
		if user != nil {
			t.Errorf("expected nil pointer when attempting to load non-existent user from context")
		}
	})

	t.Run("With no context", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContext(r)
		if err == nil {
			t.Errorf("expected error when attempting to call with no context")
		}
		if err != routes.ErrMissingUserContext {
			t.Errorf("error does not match the expected type ErrMissingUserContext")
		}
		if user != nil {
			t.Errorf("expected nil pointer when attempting to call with no context")
		}
	})
}
func TestUserContextHelpers_LoadUserFromContextWithRole(t *testing.T) {
	t.Run("Admin guard pass", func(t *testing.T) {
		payload := &dtos.JwtPayload{
			Id:   "admin",
			Role: "admin",
		}
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContextWithRole(r.WithContext(
			context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, payload),
		), constants.AdminRole)
		if err != nil {
			t.Errorf(err.Error())
		}
		if user == nil {
			t.Errorf("failed to load mock user")
		}
	})

	t.Run("User doesn't exist", func(t *testing.T) {
		payload := &dtos.JwtPayload{
			Id:   "test",
			Role: "user",
		}
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContextWithRole(r.WithContext(
			context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, payload),
		), constants.AdminRole)
		if err == nil {
			t.Errorf("expected error when attempting to load non-existent user from context")
		}
		if user != nil {
			t.Errorf("expected nil pointer when attempting to load non-existent user from context")
		}
	})

	t.Run("Admin guard block", func(t *testing.T) {
		payload := &dtos.JwtPayload{
			Id:   "user",
			Role: "user",
		}
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContextWithRole(r.WithContext(
			context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, payload),
		), constants.AdminRole)
		if err == nil {
			t.Errorf("expected error when attempting to load user with incorrect role")
		}
		if user != nil {
			t.Errorf("expected nil pointer when attempting to load user with incorrect role")
		}
	})

	t.Run("With no context", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)

		user, err := userContextHelper.LoadUserFromContextWithRole(r, constants.AdminRole)
		if err == nil {
			t.Errorf("expected error when attempting to call with no context")
		}
		if user != nil {
			t.Errorf("expected nil pointer when attempting to call with no context")
		}
	})
}
