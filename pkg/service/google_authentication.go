package service

import (
	"database/sql"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils/google"
)

// GoogleAuthenicationService for signing up and logging in users using googles authentication
type GoogleAuthenicationService interface {
	ValidateGoogleSignUp(dto *dtos.GoogleSignUpRequest) (*dtos.LoginSuccess, error)
	ValidateGoogleSignIn(string) (*dtos.LoginSuccess, error)
}

type GoogleAuthenticationConfiguration struct {
	ClientId string
}

type GoogleJsonWebTokenAuthenticationService struct {
	logger     logging.Logger
	config     *GoogleAuthenticationConfiguration
	jwtService JsonWebTokenService
	userRepo   repository.UserRepository
}

func NewGoogleAuthenticationService(config *GoogleAuthenticationConfiguration, userRepo repository.UserRepository, jwtService JsonWebTokenService, lw logging.LogWriter) GoogleAuthenicationService {
	return &GoogleJsonWebTokenAuthenticationService{
		logger:     logging.NewContextLogger(lw, "GoogleAuthenticationService"),
		config:     config,
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

// validates google ID token and logs in google user
func (svc *GoogleJsonWebTokenAuthenticationService) ValidateGoogleSignIn(idToken string) (*dtos.LoginSuccess, error) {
	if svc.config.ClientId == "" {
		svc.logger.Error(ErrGoogleClietIdNotSet, "google client id not configured")
		return nil, ErrGoogleClietIdNotSet
	}

	payload, err := google.NewValidator().ValidateToken(idToken, svc.config.ClientId)

	if err != nil {
		svc.logger.Error(ErrInvalidGoogleToken, "google id token validation failed")
		return nil, ErrInvalidGoogleToken
	}

	claims := payload.GetClaims()

	svc.logger.Debugf("claims in token: id: %s, email: %s, email_verified: %v, picture: %v", claims.Id, claims.Email, claims.EmailVerified, claims.Picture)

	existingUser, err := svc.userRepo.GetUserByEmail(claims.Email)

	if err != nil {
		svc.logger.Errorf(err, "unable to find user with email: %s", claims.Email)
		return nil, ErrUserNotFound
	}

	if existingUser.GoogleId.String != claims.Id {
		return nil, ErrUserNotFound
	}

	svc.logger.Infof("successfully verified google user %s", claims.Email)

	token, err := svc.jwtService.Sign(JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	})

	if err != nil {
		svc.logger.Error(err, "unable to sign and generate token")
		return nil, err
	}

	return &dtos.LoginSuccess{
		User: dtos.LoginUser{
			ID:        existingUser.ID,
			Username:  existingUser.Username,
			FirstName: existingUser.FirstName.String,
			LastName:  existingUser.LastName.String,
			Role:      existingUser.Role,
		},
		AccessToken: *token,
	}, nil
}

// validates google ID token and registeres google user
func (svc *GoogleJsonWebTokenAuthenticationService) ValidateGoogleSignUp(dto *dtos.GoogleSignUpRequest) (*dtos.LoginSuccess, error) {

	if svc.config.ClientId == "" {
		svc.logger.Error(ErrGoogleClietIdNotSet, "google client id not configured")
		return nil, ErrGoogleClietIdNotSet
	}

	payload, err := google.NewValidator().ValidateToken(dto.IdToken, svc.config.ClientId)

	if err != nil {
		svc.logger.Error(ErrInvalidGoogleToken, "google id token validation failed")
		return nil, ErrInvalidGoogleToken
	}

	claims := payload.GetClaims()

	svc.logger.Debugf("claims in token: id: %s, email: %s, email_verified: %v, picture: %v", claims.Id, claims.Email, claims.EmailVerified, claims.Picture)

	existingUser, _ := svc.userRepo.GetUserByEmail(claims.Email)

	if existingUser != nil {
		svc.logger.Errorf(ErrUserAlreadyExists, "user exists with provided google email- %s", claims.Email)
		return nil, ErrUserAlreadyExists
	}

	var fname, lname = dto.FirstName, dto.LastName

	if fname == "" {
		fname = claims.FirstName
	}

	if lname == "" {
		lname = claims.LastName
	}

	model := &models.UserModel{
		Email:     claims.Email,
		Password:  "",
		FirstName: sql.NullString{String: fname, Valid: true},
		LastName:  sql.NullString{String: lname, Valid: true},
		Username:  dto.Username,
		GoogleId:  sql.NullString{String: claims.Id, Valid: true},
		AvatarUrl: sql.NullString{String: claims.Picture, Valid: true},
	}

	err = svc.userRepo.InsertUser(model)

	if err != nil {
		svc.logger.Errorf(err, "unable to insert user - %v", err)
		return nil, err
	}

	token, err := svc.jwtService.Sign(JwtPayload{
		Id:   model.ID,
		Role: model.Role,
	})

	if err != nil {
		svc.logger.Error(err, "unable to sign and generate token")
		return nil, err
	}

	return &dtos.LoginSuccess{
		User: dtos.LoginUser{
			ID:        model.ID,
			Username:  model.Username,
			FirstName: model.FirstName.String,
			LastName:  model.LastName.String,
			Role:      model.Role,
		},
		AccessToken: *token,
	}, nil
}
