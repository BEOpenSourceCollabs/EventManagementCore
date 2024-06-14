package models

import (
	"database/sql"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/dtos"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/utils"
)

// UserModel represents the user data stored in the database.
type UserModel struct {
	Model
	Username  string         `db:"username"`
	Email     string         `db:"email"`
	Password  string         `db:"password" json:"-"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	BirthDate sql.NullTime   `db:"birth_date"`
	Role      constants.Role `db:"role"`
	Verified  bool           `db:"verified"`
	About     sql.NullString `db:"about"`
}

// BeforeCreate overrides model lifecycle hook, hashes the users password before proceeding.
func (m *UserModel) BeforeCreate() error {
	logger.AppLogger.InfoF("UserModel", "overridden lifecycle BeforeCreate() - hashing password")
	hash, err := utils.HashPassword(m.Password)
	if err != nil {
		return err
	}
	m.Password = hash
	return nil
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *UserModel) BeforeUpdated() error {
	logger.AppLogger.InfoF("UserModel", "overridden lifecycle BeforeUpdated() - updating updated_at")
	m.UpdatedAt = time.Now()
	return nil
}

func (m *UserModel) ToDto() dtos.LoginUser {
	return dtos.LoginUser{
		ID:        m.ID,
		FirstName: m.FirstName.String,
		LastName:  m.LastName.String,
		Username:  m.Username,
		Role:      m.Role,
	}
}

func (m *UserModel) UpdateFrom(payload dtos.CreateOrUpdateUser) {
	if len(payload.Email) > 0 {
		m.Email = payload.Email
	}
	if len(payload.Password) > 0 {
		m.Password = payload.Password
	}
	if len(payload.FirstName) > 0 {
		m.FirstName = sql.NullString{
			String: payload.FirstName,
			Valid:  true,
		}
	}
	if len(payload.LastName) > 0 {
		m.LastName = sql.NullString{
			String: payload.LastName,
			Valid:  true,
		}
	}
	if len(payload.Username) > 0 {
		m.Username = payload.Username
	}
	if len(payload.Role) > 0 {
		m.Role = payload.Role
	}
	if payload.Verified != nil {
		m.Verified = *payload.Verified
	}
}
