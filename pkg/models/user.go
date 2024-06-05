package models

import "time"

// UserModel represents the user data stored in the database.
type UserModel struct {
	ID        string     `db:"id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	BirthDate *time.Time `db:"birth_date"`
	Role      string     `db:"role"`
	Verified  bool       `db:"verified"`
	About     string     `db:"about"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}
