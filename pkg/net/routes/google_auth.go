package routes

import (
	"errors"
	"net/http"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/service"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

type googleAuthenticationRoutes struct {
	gAuthService service.GoogleAuthenicationService
	logger       logging.Logger
}

// NewGoogleAuthenticationRoutes creates routes using GoogleAuthenicationService and mounts them to the provided router.
func NewGoogleAuthenticationRoutes(router net.AppRouter, gAuthService service.GoogleAuthenicationService, lw logging.LogWriter) *googleAuthenticationRoutes {
	routes := &googleAuthenticationRoutes{
		gAuthService: gAuthService,
		logger:       logging.NewContextLogger(lw, "GoogleAuthenticationRoutes"),
	}

	router.Post("/api/auth/google/signup", http.HandlerFunc(routes.HandleGoogleSignUp))
	router.Post("/api/auth/google/signin", http.HandlerFunc(routes.HandleGoogleSignIn))

	// Add basic preflight handlers
	router.Options("/api/auth/google/signup", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	router.Options("/api/auth/google/signin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return routes
}

// handles google sign up
func (authRouter *googleAuthenticationRoutes) HandleGoogleSignUp(w http.ResponseWriter, r *http.Request) {
	gsignUpReq := &dtos.GoogleSignUpRequest{}

	//read json into registerDto
	if err := utils.ReadJson(w, r, gsignUpReq); err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	// custom validation
	if validationErrs := gsignUpReq.Validate(); len(validationErrs) > 0 {
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, validationErrs)
		return
	}

	result, err := authRouter.gAuthService.ValidateGoogleSignUp(gsignUpReq)

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

// handles google sign in
func (authRouter *googleAuthenticationRoutes) HandleGoogleSignIn(w http.ResponseWriter, r *http.Request) {

	gsignInReq := &dtos.GoogleSignInRequest{}

	//read json into registerDto
	if err := utils.ReadJson(w, r, gsignInReq); err != nil {
		utils.WriteRequestPayloadError(err, w)
		return
	}

	// custom validation
	if validationErrs := gsignInReq.Validate(); len(validationErrs) > 0 {
		utils.WriteErrorJsonResponse(w, constants.ErrorCodes.BadRequest, http.StatusBadRequest, validationErrs)
		return
	}

	result, err := authRouter.gAuthService.ValidateGoogleSignIn(gsignInReq.IdToken)

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
