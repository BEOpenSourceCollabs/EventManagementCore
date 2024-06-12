package models

import (
	"database/sql"
	"time"
)

// UserModel represents the user data stored in the database.
type UserModel struct {
	ID        string         `db:"id"`
	GoogleId  sql.NullString `db:"google_id" json:"-"`
	AvatarUrl sql.NullString `db:"avatar_url"`
	Username  string         `db:"username"`
	Email     string         `db:"email"`
	Password  string         `db:"password" json:"-"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	BirthDate sql.NullTime   `db:"birth_date"`
	Role      string         `db:"role"`
	Verified  bool           `db:"verified"`
	About     sql.NullString `db:"about"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}
