package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/middleware"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type jwtAuthRoutes struct {
	authService service.AuthenticationService
	logger      logging.Logger
}

// NewJsonWebTokenAuthenticationRoutes creates routes using AuthenticationService and JsonWebTokenService then mounts them to the provided router.
func NewJsonWebTokenAuthenticationRoutes(router net.AppRouter, authService service.AuthenticationService, jwtService *service.JsonWebTokenService, lw logging.LogWriter) *jwtAuthRoutes {
	routes := &jwtAuthRoutes{
		authService: authService,
		logger:      logging.NewContextLogger(lw, "AuthRoutes"),
	}

	protectMiddleware := middleware.JWTBearerMiddleware{
		Logger:     logging.NewContextLogger(lw, "UserRoutes.JWTBearerMiddleware"),
		JWTService: *jwtService,
	}

	router.Post("/api/auth/login", http.HandlerFunc(routes.HandleLogin))
	router.Post("/api/auth/register", http.HandlerFunc(routes.HandleSignUp))
	router.Get("/api/auth/check", protectMiddleware.BeforeNext(http.HandlerFunc(routes.HandleCheck)))
	router.Get("/api/auth/refresh", http.HandlerFunc(routes.HandleRefresh))

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

	router.Options("/api/auth/refresh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return routes
}

// HandleLogin user login
func (authRouter *jwtAuthRoutes) HandleLogin(w http.ResponseWriter, r *http.Request) {
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

	//attach refresh token cookie to response
	err = authRouter.authService.AttachRefreshTokenCookie(w, data.User.ID)

	if err != nil {
		authRouter.logger.Error(err, "error attaching refresh token cookie")
		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	//return access token
	utils.WriteSuccessJsonResponse(w, http.StatusOK, data)
}

// HandleSignUp user sign up
func (authRouter *jwtAuthRoutes) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	//initialize register struct
	registerDto := &dtos.Register{}

	//read json into registerDto
	if err := utils.ReadJson(w, r, registerDto); err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	// custom validation
	if validationErrs := registerDto.Validate(); len(validationErrs) > 0 {
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, validationErrs)
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

// HandleCheck checks if access token is valid
func (authRouter *jwtAuthRoutes) HandleCheck(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(service.USER_CONTEXT_KEY).(*service.JwtPayload)

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

// Handles validating refresh token & providing new access token
func (authRouter *jwtAuthRoutes) HandleRefresh(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie(constants.REFRESH_TOKEN_COOKIE)

	if err != nil {
		authRouter.logger.Error(err, "error parsing refresh token cookie")

		if err == http.ErrNoCookie {
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthNoRefreshTokenCookie, http.StatusUnauthorized, []string{"No refresh token provided"})
			return
		}

		utils.WriteInternalErrorJsonResponse(w)
		return
	}

	accessToken, err := authRouter.authService.ValidateRefresh(cookie.Value)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.AuthInvalidRefreshToken, http.StatusUnauthorized, nil)
			return
		case errors.Is(err, service.ErrUserNotFound):
			utils.WriteErrorJsonResponse(w, constants.ErrorCodes.NotFound, http.StatusNotFound, []string{err.Error()})
			return
		default:
			utils.WriteInternalErrorJsonResponse(w)
			return
		}
	}

	utils.WriteSuccessJsonResponse(w, http.StatusOK, accessToken)

}
