package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type authRoutes struct {
	config      *service.AuthServiceConfiguration
	authService service.IAuthService
}

func NewAuthRoutes(router net.AppRouter, authService service.IAuthService, config *service.AuthServiceConfiguration) *authRoutes {
	routes := &authRoutes{
		authService: authService,
		config:      config,
	}

	protectMiddleware := middleware.JWTBearerMiddleware{
		Secret: config.Secret,
	}

	router.Post("/api/auth/login", http.HandlerFunc(routes.HandleLogin))
	router.Post("/api/auth/register", http.HandlerFunc(routes.HandleSignUp))
	router.Get("/api/auth/check", protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleCheck)))

	// Add basic preflight handlers
	router.Options("/api/auth/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	router.Options("/api/auth/register", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	router.Options("/api/auth/check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	router.Post("/api/auth/google/signin", http.HandlerFunc(routes.HandleGoogleSignIn))
	router.Post("/api/auth/google/signup", http.HandlerFunc(routes.HandleGoogleSignUp))

	return routes
}

// handles user login
func (authRouter *authRoutes) HandleLogin(w http.ResponseWriter, r *http.Request) {
	//initialize login struct
	loginDto := &dtos.Login{}

	//read json into loginDto
	err := utils.ReadJson(w, r, loginDto)

	if err != nil {
		utils.WriteRequestPayloadError(err, w)
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
		utils.WriteRequestPayloadError(err, w)
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

	user := r.Context().Value(constants.USER_CONTEXT_KEY).(*dtos.JwtPayload)

	result, err := authRouter.authService.CheckUser(user.Id)

	if err != nil {

		if errors.Is(err, service.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{fmt.Sprintf("user with id %s does not exist", user.Id)})
			return
		}

		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, result)
}

// handles google sign in
func (authRouter *authRoutes) HandleGoogleSignIn(w http.ResponseWriter, r *http.Request) {

	gsignInReq := &dtos.GoogleSignInRequest{}

	//read json into registerDto
	err := utils.ReadJson(w, r, gsignInReq)

	if err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	err = utils.Validator.Struct(gsignInReq)

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

	result, err := authRouter.authService.ValidateGoogleSignIn(gsignInReq.IdToken)

	if err != nil {

		if errors.Is(err, service.ErrInvalidGoogleToken) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidAuthToken, http.StatusUnauthorized, []string{"google id token is invalid"})
			return
		}

		if errors.Is(err, service.ErrUserNotFound) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{"user does not exist with provided google account"})
			return
		}

		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, result)

}

// handles google sign up
func (authRouter *authRoutes) HandleGoogleSignUp(w http.ResponseWriter, r *http.Request) {
	gsignUpReq := &dtos.GoogleSignUpRequest{}

	//read json into registerDto
	err := utils.ReadJson(w, r, gsignUpReq)

	if err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	err = utils.Validator.Struct(gsignUpReq)

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

	result, err := authRouter.authService.ValidateGoogleSignUp(gsignUpReq)

	if err != nil {

		if errors.Is(err, service.ErrInvalidGoogleToken) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidAuthToken, http.StatusUnauthorized, []string{"google id token is invalid"})
			return
		}

		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, []string{"user account with email already exist"})
			return
		}

		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	utils.WriteSuccessJsonResponse(w, http.StatusCreated, result)

}
