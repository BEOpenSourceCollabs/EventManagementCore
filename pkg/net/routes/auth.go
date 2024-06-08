package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type authRoutes struct {
	logger      *log.Logger
	authService *service.AuthService
}

func NewAuthRoutes(router net.AppRouter, authService *service.AuthService, logger *log.Logger) *authRoutes {

	routes := &authRoutes{
		authService: authService,
		logger:      logger,
	}

	router.Post("/api/auth/login", http.HandlerFunc(routes.HandleLogin))
	router.Post("/api/auth/register", http.HandlerFunc(routes.HandleSignUp))
	router.Get("/api/auth/check", http.HandlerFunc(routes.HandleCheck))

	return routes
}

// handles user login
func (authRouter *authRoutes) HandleLogin(w http.ResponseWriter, r *http.Request) {
	//initialize login struct
	loginDto := &dtos.Login{}

	//read json into loginDto
	err := utils.ReadJson(w, r, loginDto)

	if err != nil {
		utils.HandleCommonJsonError(err, w)
		return
	}

	//validate login
	data, err := authRouter.authService.ValidateSignIn(loginDto)

	//handle errors
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidCredentials, http.StatusUnauthorized, nil)
			return
		default:
			utils.WriteInternalErrorJsonResponse(w)
			return
		}
	}

	//return access token
	utils.WriteSuccessJsonResponse(w, http.StatusOK, data)
}

// handles user sign up
func (authRouter *authRoutes) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	//initialize register struct
	registerDto := &dtos.Register{}

	//read json into registerDto
	err := utils.ReadJson(w, r, registerDto)

	if err != nil {
		utils.HandleCommonJsonError(err, w)
		return
	}

	//custom validation
	err = utils.Validator.Struct(registerDto)

	if err != nil {
		var validationErrors = []utils.ValidationErrorResponse{}
		for _, err := range err.(validator.ValidationErrors) {
			validationErr := utils.ValidationErrorResponse{
				Field:   err.Field(),
				Rule:    err.ActualTag(),
				Value:   err.Value(),
				Message: utils.HumanFriendlyErrorMessage(err.ActualTag(), err.Param()),
			}
			validationErrors = append(validationErrors, validationErr)
		}
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, validationErrors)
		return
	}

	result, err := authRouter.authService.ValidateSignUp(registerDto)

	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{fmt.Sprintf("User already exists with email %s", registerDto.Email)})
			return
		}

		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusCreated, result)
}

// checks if access token is valid
func (authRouter *authRoutes) HandleCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteErrorJsonResponse(w, "Not implemented", http.StatusInternalServerError, nil)
}
