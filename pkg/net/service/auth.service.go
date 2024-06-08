package service

import (
	"errors"
	"log"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists with given email")
)

type AuthService struct {
	config   *config.Configuration
	logger   *log.Logger
	userRepo repository.UserRepository
}

func NewAuthService(config *config.Configuration, logger *log.Logger, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		config:   config,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (svc *AuthService) ValidateSignIn(dto *dtos.Login) (*dtos.LoginSuccess, error) {

	//check if user exists with provided email
	existingUser, err := svc.userRepo.GetUserByEmail(dto.Email)

	//if not return invalid credentials error.
	if err != nil {
		svc.logger.Printf("user with email %s not found in db", dto.Email)
		return nil, ErrInvalidCredentials
	}

	// check if provided password and password from db match
	//TODO: use hashed passwords for now checking plain text passwords.
	if existingUser.Password != dto.Password {
		svc.logger.Printf("password didn't match for user with email %s", dto.Email)
		return nil, ErrInvalidCredentials
	}

	//TODO: create a proper JWT token
	return &dtos.LoginSuccess{
		User: dtos.LoginUser{
			ID:        existingUser.ID,
			Username:  existingUser.Username,
			FirstName: existingUser.FirstName,
			LastName:  existingUser.LastName,
			Role:      existingUser.Role,
		},
		AccessToken: "token",
	}, nil
}

func (svc *AuthService) ValidateSignUp(dto *dtos.Register) (string, error) {

	//check if user exists with provided email
	existingUser, _ := svc.userRepo.GetUserByEmail(dto.Email)

	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	model := &models.UserModel{
		Email:     dto.Email,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Username:  dto.Username,
	}

	err := svc.userRepo.InsertUser(model)

	if err != nil {
		return "", err
	}

	return "successfully signed up", nil
}
