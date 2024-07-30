package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// CategoryRepository represents the interface for category-related database operations.
type CategoryRepository interface {
	CreateCategory(category *models.CategoryModel) error
	GetCategoryByID(id string) (*models.CategoryModel, error)
	UpdateCategory(category *models.CategoryModel) error
	DeleteCategory(id string) error
}

type sqlCategoryRepository struct {
	database *sql.DB
}

// NewSQLCategoryRepository creates and returns a new sql flavoured CategoryRepository instance.
func NewSQLCategoryRepository(database *sql.DB) CategoryRepository {
	return &sqlCategoryRepository{database: database}
}

// inserts a new category into the database
func (r *sqlCategoryRepository) CreateCategory(category *models.CategoryModel) error {
	category.BeforeCreate()

	query := `INSERT INTO public.categories (name) VALUES ($1) RETURNING id`

	err := r.database.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	category.AfterCreate()
	return nil
}

// retrieves a category from the database by its unique ID, if it exists
func (r *sqlCategoryRepository) GetCategoryByID(id string) (*models.CategoryModel, error) {
	query := `SELECT
				id,
				name

			FROM public.categories WHERE id = $1`

	category := &models.CategoryModel{}
	err := r.database.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidCategoryId
	}

	return category, nil
}

// updates a category in database
func (r *sqlCategoryRepository) UpdateCategory(category *models.CategoryModel) error {
	category.BeforeUpdate()

	query := `UPDATE public.categories SET name = $1 WHERE id = $2`
	rs, err := r.database.Exec(
		query,
		category.Name,
		category.ID,
	)
	if err != nil {
		return err
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrCategoryNotFound
	}

	category.AfterUpdate()
	return nil
}

// deletes a category from the database.
func (r *sqlCategoryRepository) DeleteCategory(id string) error {
	query := `DELETE FROM public.categories WHERE id = $1`

	rs, err := r.database.Exec(query, id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrCategoryNotFound
	}

	return nil
}

// errors
var (
	ErrCategoryNotFound  = errors.New("category not found")
	ErrInvalidCategoryId = errors.New("invalid category id")
)
