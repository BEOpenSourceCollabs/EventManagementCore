package net

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

// UserContextHelpers provides pluggable helpers to route structures that use, user context within requests.
type UserContextHelpers struct {
	R *repository.UserRepository
}

// LoadUserFromContext helper that attempts to read the http.Request's user context key or returns an error if it was not found.
// Returns the loaded user if found.
func (h UserContextHelpers) LoadUserFromContext(r *http.Request) (*models.UserModel, error) {
	userContext, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*dtos.JwtPayload)
	if !ok {
		return nil, ErrMissingUserContext
	}
	return (*h.R).GetUserByID(userContext.Id)
}

// LoadUserFromContextWithRole helper that attempts to read the http.Request's user context key or returns an error if it was not found.
// Returns the loaded user if found and has the role specified in the parameters.
// This helper can be used as a gaurd to protect routes being accessed by users without the specified role.
func (h UserContextHelpers) LoadUserFromContextWithRole(r *http.Request, role constants.Role) (*models.UserModel, error) {
	user, err := h.LoadUserFromContext(r)
	if err != nil {
		return nil, err
	}
	if user.Role != role {
		return nil, fmt.Errorf("user is missing the required role '%v'", role)
	}
	return user, nil
}

var (
	ErrMissingUserContext = errors.New("no user context provided") // ErrMissingUserContext is returned when no context is found while attempting to load user from http.Requests context.
)
