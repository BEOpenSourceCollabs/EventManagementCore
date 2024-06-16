package mock

import "github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"

type UserRepository struct {
	CreateUserFn     func(user *models.UserModel) error
	GetUserByIDFn    func(id string) (*models.UserModel, error)
	UpdateUserFn     func(user *models.UserModel) error
	DeleteUserFn     func(id string) error
	GetUserByEmailFn func(email string) (*models.UserModel, error)
	InsertUserFn     func(user *models.UserModel) error
}

func (u UserRepository) CreateUser(user *models.UserModel) error {
	if u.CreateUserFn != nil {
		return u.CreateUserFn(user)
	}
	return nil
}

func (u UserRepository) GetUserByID(id string) (*models.UserModel, error) {
	if u.GetUserByIDFn != nil {
		return u.GetUserByIDFn(id)
	}
	return nil, nil
}

func (u UserRepository) UpdateUser(user *models.UserModel) error {
	if u.UpdateUserFn != nil {
		return u.UpdateUserFn(user)
	}
	return nil
}

func (u UserRepository) DeleteUser(id string) error {
	if u.DeleteUserFn != nil {
		return u.DeleteUserFn(id)
	}
	return nil
}

func (u UserRepository) GetUserByEmail(email string) (*models.UserModel, error) {
	if u.GetUserByEmailFn != nil {
		return u.GetUserByEmailFn(email)
	}
	return nil, nil
}

func (u UserRepository) InsertUser(user *models.UserModel) error {
	if u.InsertUserFn != nil {
		return u.InsertUserFn(user)
	}
	return nil
}
