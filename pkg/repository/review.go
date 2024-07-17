package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

type ReviewRepository interface {
	CreateReview(review *models.ReviewModel) error
	GetReviewByID(id string) (*models.ReviewModel, error)
	DeleteReview(id string) error
	UpdateReview(review *models.ReviewModel) error
}

type sqlReviewRepository struct {
	database *sql.DB
}

func NewSQLReviewRepository(database *sql.DB) ReviewRepository {
	return &sqlReviewRepository{database: database}
}

func (r *sqlReviewRepository) handleError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrReviewNotFound
	}
	if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
		return ErrRepoConnErr
	}
	return err
}

// CreateReview inserts a new review into the database
func (r *sqlReviewRepository) CreateReview(review *models.ReviewModel) error {
	review.BeforeCreate()
	query := `INSERT INTO public.reviews (title, event_id, author_id, body)
        VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.database.QueryRow(query,
		review.Title,
		review.EventID,
		review.AuthorID,
		review.Body,
	).Scan(&review.ID)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	review.AfterCreate()
	return nil
}

// GetReviewByID retrieves a review from the database by its unique ID
func (r *sqlReviewRepository) GetReviewByID(id string) (*models.ReviewModel, error) {
	query := `SELECT id, title, event_id, author_id, body, created_at, updated_at FROM public.reviews WHERE id = $1`
	review := &models.ReviewModel{}
	err := r.database.QueryRow(query, id).Scan(
		&review.ID,
		&review.Title,
		&review.EventID,
		&review.AuthorID,
		&review.Body,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return nil, r.handleError(err)
	}
	return review, nil
}

// UpdateReview updates a review in the database
func (r *sqlReviewRepository) UpdateReview(review *models.ReviewModel) error {
	review.BeforeUpdate()
	query := `UPDATE public.reviews SET title = $1, event_id = $2, author_id = $3, body = $4 WHERE id = $5`
	if review.CreatedAt.Unix() == 0 {
		return fmt.Errorf("unable to update a review that was not loaded from the database")
	}
	rs, err := r.database.Exec(
		query,
		review.Title,
		review.EventID,
		review.AuthorID,
		review.Body,
		review.ID,
	)
	if err != nil {
		return r.handleError(err)
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrReviewNotFound
	}
	review.AfterUpdate()
	return nil
}

// DeleteReview deletes a review from the database
func (r *sqlReviewRepository) DeleteReview(id string) error {
	query := `DELETE FROM public.reviews WHERE id = $1`
	rs, err := r.database.Exec(query, id)
	if err != nil {
		return r.handleError(err)
	}
	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrReviewNotFound
	}
	return nil
}

var (
	ErrReviewNotFound = errors.New("review not found") // ErrReviewNotFound is returned when a review is not found in the database.
)
