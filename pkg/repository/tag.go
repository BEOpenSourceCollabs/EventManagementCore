package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// TagRepository represents the interface for tag-related database operations.
type TagRepository interface {
	CreateTag(tag *models.TagModel) error
	GetTagByID(id string) (*models.TagModel, error)
	UpdateTag(tag *models.TagModel) error
	DeleteTag(id string) error
}

type sqlTagRepository struct {
	database *sql.DB
}

// NewSQLTagRepository creates and returns a new sql flavoured TagRepository instance.
func NewSQLTagRepository(database *sql.DB) TagRepository {
	return &sqlTagRepository{database: database}
}

// inserts a new tag into the database
func (r *sqlTagRepository) CreateTag(tag *models.TagModel) error {
	tag.BeforeCreate()

	query := `INSERT INTO public.tags (name) VALUES ($1) RETURNING id`

	err := r.database.QueryRow(query, tag.Name).Scan(&tag.ID)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	tag.AfterCreate()
	return nil
}

// retrieves a tag from the database by its unique ID, if it exists
func (r *sqlTagRepository) GetTagByID(id string) (*models.TagModel, error) {
	query := `SELECT
				id,
				name

			FROM public.tags WHERE id = $1`

	tag := &models.TagModel{}
	err := r.database.QueryRow(query, id).Scan(
		&tag.ID,
		&tag.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTagNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidTagId
	}

	return tag, nil
}

// updates a tag in database
func (r *sqlTagRepository) UpdateTag(tag *models.TagModel) error {
	tag.BeforeUpdate()

	query := `UPDATE public.tags SET name = $1 WHERE id = $2`
	rs, err := r.database.Exec(
		query,
		tag.Name,
		tag.ID,
	)
	if err != nil {
		return err
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrTagNotFound
	}

	tag.AfterUpdate()
	return nil
}

// deletes a tag from the database.
func (r *sqlTagRepository) DeleteTag(id string) error {
	query := `DELETE FROM public.tags WHERE id = $1`

	rs, err := r.database.Exec(query, id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrTagNotFound
	}

	return nil
}

// errors
var (
	ErrTagNotFound  = errors.New("tag not found")
	ErrInvalidTagId = errors.New("invalid tag id")
)
