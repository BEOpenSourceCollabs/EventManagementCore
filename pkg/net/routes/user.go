package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type userRoutes struct {
	UserContextHelpers // include user context helpers
	userRepository     repository.UserRepository
}

func NewUserRoutes(router net.AppRouter, userRepository repository.UserRepository, authConfig *service.AuthServiceConfiguration) userRoutes {
	routes := userRoutes{
		/* inject dependencies */
		userRepository: userRepository,
		UserContextHelpers: UserContextHelpers{
			r: &userRepository,
		},
	}

	// mount routes to router.
	router.Post(
		"/api/users",
		middleware.ProtectMiddleware(
			http.HandlerFunc(routes.HandleCreateUser),
			authConfig.Secret,
		),
	)
	router.Get(
		"/api/users/{id}",
		middleware.ProtectMiddleware(
			http.HandlerFunc(routes.HandleGetUserById),
			authConfig.Secret,
		),
	)
	router.Put(
		"/api/users/{id}",
		middleware.ProtectMiddleware(
			http.HandlerFunc(routes.HandleUpdateUserById),
			authConfig.Secret,
		),
	)
	router.Delete(
		"/api/users/{id}",
		middleware.ProtectMiddleware(
			http.HandlerFunc(routes.HandleDeleteUserById),
			authConfig.Secret,
		),
	)
	return routes
}

func (u userRoutes) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, constants.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, err.Error())
		return
	}

	payload := &dtos.CreateUser{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, err.Error())
		return
	}

	hash, _ := utils.HashPassword(payload.Password)

	// TODO: validation
	// No validation is done while creating a user as admin.
	// To do so would involve including the authentication service as a dependency.
	user := &models.UserModel{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hash,
		FirstName: sql.NullString{
			String: payload.FirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: payload.LastName,
			Valid:  true,
		},
		Role:     payload.Role,
		Verified: payload.Verified,
	}

	if err := u.userRepository.CreateUser(user); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.InternalServerError, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, "user created")
}

func (u userRoutes) HandleGetUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, constants.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, err.Error())
		return
	}

	id := r.PathValue("id")

	user, err := u.userRepository.GetUserByID(id)

	if err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (u userRoutes) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, constants.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, err.Error())
		return
	}

	id := r.PathValue("id")

	// Load the user first
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Merge in the request payload
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// submit the changes
	if err := u.userRepository.UpdateUser(user); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

func (u userRoutes) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	// Ensure that a valid user with the "admin" role is accessing this api.
	if _, err := u.LoadUserFromContextWithRole(r, constants.AdminRole); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidScope, http.StatusUnauthorized, err.Error())
		return
	}

	id := r.PathValue("id")

	if err := u.userRepository.DeleteUser(id); err != nil {
		logger.AppLogger.Error("userRoutes", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
