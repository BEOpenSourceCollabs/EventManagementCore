package service

import (
	"database/sql"
	"errors"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logging"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

const (
	USER_CONTEXT_KEY types.ContextKey = "user"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists with given email")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidGoogleToken  = errors.New("google id token is invalid")
	ErrGoogleClietIdNotSet = errors.New("google client id not configured")
)

// AuthenticationService for signing up and logging in users.
type AuthenticationService interface {
	ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error)
	ValidateSignUp(dto *dtos.Register) (string, error)
	CheckUser(id string) (*dtos.LoginUser, error)
}

// jsonWebTokenAuthenticationService implementation of the AuthenticationService using the JsonWebTokenService.
type jsonWebTokenAuthenticationService struct {
	logger     logging.Logger
	jwtService JsonWebTokenService
	userRepo   repository.UserRepository
}

// NewJsonWebTokenAuthenticationService create a JWT flavoured AuthenticationService.
func NewJsonWebTokenAuthenticationService(userRepo repository.UserRepository, jwtService JsonWebTokenService, lw logging.LogWriter) AuthenticationService {
	return &jsonWebTokenAuthenticationService{
		logger:     logging.NewContextLogger(lw, "JsonWebTokenAuthenticationService"),
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (svc *jsonWebTokenAuthenticationService) ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error) {

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

	token, err := svc.jwtService.SignAccessToken(JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	})

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
		AccessToken: *token,
	}, nil
}

func (svc *jsonWebTokenAuthenticationService) ValidateSignUp(dto *dtos.Register) (string, error) {

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

func (svc *jsonWebTokenAuthenticationService) CheckUser(id string) (*dtos.LoginUser, error) {

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
