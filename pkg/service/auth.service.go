package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/config"
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
		svc.logger.Printf("%v", err)
		svc.logger.Printf("user with email %s not found in db", dto.Email)
		return nil, ErrInvalidCredentials
	}

	// check if provided password and password from db match
	if !utils.DoesPasswordMatch(dto.Password, existingUser.Password) {
		svc.logger.Printf("password didn't match for user with email %s", dto.Email)
		return nil, ErrInvalidCredentials
	}

	jwtPayload := &dtos.JwtPayload{
		Id:   existingUser.ID,
		Role: existingUser.Role,
	}

	token, err := utils.GenerateToken(jwtPayload, svc.config.Secret)

	if err != nil {
		svc.logger.Printf("%v", err)
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

	hashedPw, err := utils.HashPassword(dto.Password)

	if err != nil {
		return "", err
	}

	model := &models.UserModel{
		Email:     dto.Email,
		Password:  hashedPw,
		FirstName: sql.NullString{String: dto.FirstName, Valid: true},
		LastName:  sql.NullString{String: dto.LastName, Valid: true},
		Username:  dto.Username,
	}

	err = svc.userRepo.InsertUser(model)

	if err != nil {
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
