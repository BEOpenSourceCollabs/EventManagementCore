package service

import (
	"database/sql"
	"errors"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists with given email")
	ErrUserNotFound       = errors.New("user not found")
)

type IAuthService interface {
	ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error)
	ValidateSignUp(dto *dtos.Register) (string, error)
	CheckUser(id string) (*dtos.LoginUser, error)
}

type AuthServiceConfiguration struct {
	Secret string
}

type AuthService struct {
	config   *AuthServiceConfiguration
	userRepo repository.UserRepository
}

func NewAuthService(config *AuthServiceConfiguration, userRepo repository.UserRepository) IAuthService {
	return &AuthService{
		config:   config,
		userRepo: userRepo,
	}
}

func (svc *AuthService) ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error) {

	//check if user exists with provided email
	existingUser, err := svc.userRepo.GetUserByEmail(dto.Email)

	//if not return invalid credentials error.
	if err != nil {
		logger.AppLogger.ErrorF("AuthService.ValidateSignIn", "user with email %s not found in db", dto.Email)
		return nil, ErrInvalidCredentials
	}

	// check if provided password and password from db match
	if !utils.DoesPasswordMatch(dto.Password, existingUser.Password) {
		logger.AppLogger.ErrorF("AuthService.ValidateSignIn", "password didn't match for user with email %s", dto.Email)
		return nil, ErrInvalidCredentials
	}

	jwtPayload := &dtos.JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	}

	token, err := utils.GenerateToken(jwtPayload, svc.config.Secret)

	if err != nil {
		logger.AppLogger.ErrorF("AuthService.ValidateSignIn", "%v", err)
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
		logger.AppLogger.ErrorF("AuthService.ValidateSignUp", "%v", err)
		return "", err
	}

	return "successfully signed up", nil
}

func (svc *AuthService) CheckUser(id string) (*dtos.LoginUser, error) {

	existingUser, err := svc.userRepo.GetUserByID(id)

	if err != nil {
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
