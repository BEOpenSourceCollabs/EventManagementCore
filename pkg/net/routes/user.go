package routes

import (
	"errors"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type userRoutes struct {
	net.UserContextHelpers // include user context helpers
	userRepository         repository.UserRepository
}

func NewUserRoutes(router net.AppRouter, userRepository repository.UserRepository, authConfig *service.AuthServiceConfiguration) userRoutes {
	routes := userRoutes{
		/* inject dependencies */
		userRepository: userRepository,
		UserContextHelpers: net.UserContextHelpers{
			R: &userRepository,
		},
	}

	// initialize a protect middleware (factory) to wrap and protect each of the routes.
	protectMiddleware := middleware.JWTBearerMiddleware{
		Secret: authConfig.Secret,
	}

	// mount routes to router.
	router.Post(
		"/api/users",
		protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleCreateUser)),
	)
	router.Get(
		"/api/users/{id}",
		protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleGetUserById)),
	)
	router.Put(
		"/api/users/{id}",
		protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleUpdateUserById)),
	)
	router.Delete(
		"/api/users/{id}",
		protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleDeleteUserById)),
	)

	// Add basic preflight handlers
	router.Options("/api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	router.Options("/api/users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return routes
}

func (u userRoutes) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, types.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrRepoConnErr) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		} else {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, []string{err.Error()})
		}
		return
	}

	payload := dtos.CreateOrUpdateUser{}
	if err := utils.ReadJson(w, r, &payload); err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	// custom validation
	if validationErrs := payload.Validate(); len(validationErrs) > 0 {
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, validationErrs)
		return
	}

	user := &models.UserModel{}
	user.UpdateFrom(payload)

	if err := u.userRepository.CreateUser(user); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, user)
}

func (u userRoutes) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, types.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrRepoConnErr) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		} else {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, []string{err.Error()})
		}
		return
	}

	id := r.PathValue("id")

	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{err.Error()})
			return
		}
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{err.Error()})
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, user)
}

func (u userRoutes) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, types.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrRepoConnErr) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		} else {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, []string{err.Error()})
		}
		return
	}

	id := r.PathValue("id")

	// Load the user first
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{err.Error()})
			return
		}
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	payload := dtos.CreateOrUpdateUser{}
	if err := utils.ReadJson(w, r, &payload); err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	// updates the user from the payload
	user.UpdateFrom(payload)

	// submit the changes
	if err := u.userRepository.UpdateUser(user); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{err.Error()})
			return
		}
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, user)
}

func (u userRoutes) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, types.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrRepoConnErr) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		} else {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, []string{err.Error()})
		}
		return
	}

	id := r.PathValue("id")

	if err := u.userRepository.DeleteUser(id); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{err.Error()})
			return
		}
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, []string{err.Error()})
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, nil)
}
