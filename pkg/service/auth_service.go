package service

import (
	"database/sql"
	"errors"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils/google"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists with given email")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidGoogleToken  = errors.New("google id token is invalid")
	ErrGoogleClietIdNotSet = errors.New("google client id not configured")
)

type IAuthService interface {
	ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error)
	ValidateSignUp(dto *dtos.Register) (string, error)
	CheckUser(id string) (*dtos.LoginUser, error)
	ValidateGoogleSignUp(dto *dtos.GoogleSignUpRequest) (*dtos.LoginSuccess, error)
	ValidateGoogleSignIn(string) (*dtos.LoginSuccess, error)
}

type AuthServiceConfiguration struct {
	Secret         string
	GoogleClientId string
}

type AuthService struct {
	logger   logging.Logger
	config   *AuthServiceConfiguration
	userRepo repository.UserRepository
}

func NewAuthService(config *AuthServiceConfiguration, userRepo repository.UserRepository, lw logging.LogWriter) IAuthService {
	return &AuthService{
		logger:   logging.NewContextLogger(lw, "AuthService"),
		config:   config,
		userRepo: userRepo,
	}
}

func (svc *AuthService) ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error) {

	//check if user exists with provided email
	existingUser, err := svc.userRepo.GetUserByEmail(dto.Email)

	//if not return invalid credentials error.
	if err != nil {
		svc.logger.Warnf("user with email %s not found in db", dto.Email)
		return nil, ErrInvalidCredentials
	}

	// check if provided password and password from db match
	if !utils.DoesPasswordMatch(dto.Password, existingUser.Password) {
		svc.logger.Warnf("password didn't match for user with email %s", dto.Email)
		return nil, ErrInvalidCredentials
	}

	jwtPayload := &dtos.JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	}

	token, err := jwtPayload.Sign(svc.config.Secret)

	if err != nil {
		svc.logger.Error(err, "error signing token")
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
		AccessToken: token,
	}, nil
}

func (svc *AuthService) ValidateSignUp(dto *dtos.Register) (string, error) {

	//check if user exists with provided email
	existingUser, _ := svc.userRepo.GetUserByEmail(dto.Email)

	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	// NOTE: moved password hashing to model BeforeCreate lifecycle.

	model := &models.UserModel{
		Email:     dto.Email,
		Password:  dto.Password, // hashedPw,
		FirstName: sql.NullString{String: dto.FirstName, Valid: true},
		LastName:  sql.NullString{String: dto.LastName, Valid: true},
		Username:  dto.Username,
	}

	err := svc.userRepo.InsertUser(model)

	if err != nil {
		svc.logger.Error(err, "error inserting user")
		return "", err
	}

	return "successfully signed up", nil
}

func (svc *AuthService) CheckUser(id string) (*dtos.LoginUser, error) {

	existingUser, err := svc.userRepo.GetUserByID(id)

	if err != nil {
		svc.logger.Errorf(err, "unable to find user with id: %s", id)
		return nil, ErrUserNotFound
	}

	return &dtos.LoginUser{
		ID:        existingUser.ID,
		Username:  existingUser.Username,
		FirstName: existingUser.FirstName.String,
		LastName:  existingUser.LastName.String,
		Role:      existingUser.Role,
	}, nil

}

// validates google ID token and logs in google user
func (svc *AuthService) ValidateGoogleSignIn(idToken string) (*dtos.LoginSuccess, error) {

	if svc.config.GoogleClientId == "" {
		svc.logger.Error(ErrGoogleClietIdNotSet, "google client id not configured")
		return nil, ErrGoogleClietIdNotSet
	}

	payload, err := google.NewValidator().ValidateToken(idToken, svc.config.GoogleClientId)

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

	jwtPayload := &dtos.JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	}

	token, err := jwtPayload.Sign(svc.config.Secret)

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
		AccessToken: token,
	}, nil
}

// validates google ID token and registeres google user
func (svc *AuthService) ValidateGoogleSignUp(dto *dtos.GoogleSignUpRequest) (*dtos.LoginSuccess, error) {

	if svc.config.GoogleClientId == "" {
		svc.logger.Error(ErrGoogleClietIdNotSet, "google client id not configured")
		return nil, ErrGoogleClietIdNotSet
	}

	payload, err := google.NewValidator().ValidateToken(dto.IdToken, svc.config.GoogleClientId)

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

	jwtPayload := &dtos.JwtPayload{
		Id:   model.ID,
		Role: model.Role,
	}

	token, err := jwtPayload.Sign(svc.config.Secret)

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
		AccessToken: token,
	}, nil
}
