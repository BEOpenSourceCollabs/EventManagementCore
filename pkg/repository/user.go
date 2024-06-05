package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// UserRepository represents the interface for user-related database operations.
type UserRepository interface {
	CreateUser(user *models.UserModel) error
	GetUserByID(id string) (*models.UserModel, error)
	GetUserByUsername(username string) (*models.UserModel, error)
}

type sqlUserRepository struct {
	database *sql.DB
}

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(database *sql.DB) UserRepository {
	return &sqlUserRepository{database: database}
}

// CreateUser inserts a new user into the database.
func (r *sqlUserRepository) CreateUser(user *models.UserModel) error {
	query := `
		INSERT INTO public.users (username, email, password, first_name, last_name, birth_date, role, verified, about)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	err := r.database.QueryRow(query, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.BirthDate, user.Role, user.Verified, user.About).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user from the database by its unique ID.
func (r *sqlUserRepository) GetUserByID(id string) (*models.UserModel, error) {
	query := `SELECT * FROM public.users WHERE id = $1`

	user := &models.UserModel{}
	err := r.database.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.Role,
		&user.Verified,
		&user.About,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user from the database by its unique username.
func (r *sqlUserRepository) GetUserByUsername(username string) (*models.UserModel, error) {
	query := `SELECT * FROM public.users WHERE username = $1`

	user := &models.UserModel{}
	err := r.database.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.BirthDate,
		&user.Role,
		&user.Verified,
		&user.About,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

var (
	// ErrUserNotFound is returned when a user is not found in the database.
	ErrUserNotFound = errors.New("user not found")
)
