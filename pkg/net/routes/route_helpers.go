package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

type UserContextHelpers struct {
	r *repository.UserRepository
}

// LoadUserFromContext helper that attempts to read the http.Request's user context key or returns an error if it was not found.
func (h UserContextHelpers) LoadUserFromContext(r *http.Request) (*models.UserModel, error) {
	userContext, ok := r.Context().Value(service.USER_CONTEXT_KEY).(*dtos.JwtPayload)
	if !ok {
		return nil, ErrMissingUserContext
	}

	user, err := (*h.r).GetUserByID(userContext.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h UserContextHelpers) LoadUserFromContextWithRole(r *http.Request, role types.Role) (*models.UserModel, error) {
	userContext, ok := r.Context().Value(service.USER_CONTEXT_KEY).(*dtos.JwtPayload)
	if !ok {
		return nil, ErrMissingUserContext
	}

	user, err := (*h.r).GetUserByID(userContext.Id)
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
