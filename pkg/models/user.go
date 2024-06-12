package models

import (
	"database/sql"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/net/constants"
)

// UserModel represents the user data stored in the database.
type UserModel struct {
	ID        string         `db:"id"`
	Username  string         `db:"username"`
	Email     string         `db:"email"`
	Password  string         `db:"password" json:"-"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	BirthDate sql.NullTime   `db:"birth_date"`
	Role      constants.Role `db:"role"`
	Verified  bool           `db:"verified"`
	About     sql.NullString `db:"about"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
