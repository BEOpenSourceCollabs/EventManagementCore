package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// ReviewRepository represents the interface for review-related database operations.
type ReviewRepository interface {
	CreateReview(review *models.ReviewModel) error
	GetReviewByID(id string) (*models.ReviewModel, error)
	GetReviewsByEventId(event_id string) ([]models.ReviewModel, error)
	UpdateReview(review *models.ReviewModel) error
	DeleteReview(id string) error
}

type sqlReviewRepository struct {
	database *sql.DB
}

// NewSQLReviewRepository creates and returns a new sql flavoured ReviewRepository instance.
func NewSQLReviewRepository(database *sql.DB) ReviewRepository {
	return &sqlReviewRepository{database: database}
}

// inserts a new review into the database
func (r *sqlReviewRepository) CreateReview(review *models.ReviewModel) error {
	review.BeforeCreate()

	query := `INSERT INTO public.reviews (title, event_id, author_id, body)
		VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.database.QueryRow(query, review.Title, review.EventId, review.AuthorId, review.Body).Scan(&review.ID)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	review.AfterCreate()
	return nil
}

// retrieves a review from the database by its unique ID if it exists
func (r *sqlReviewRepository) GetReviewByID(id string) (*models.ReviewModel, error) {
	query := `SELECT
				id,
				title,
				event_id,
				author_id,
				body,
				created_at,
				updated_at

			FROM public.reviews WHERE id = $1`

	review := &models.ReviewModel{}
	err := r.database.QueryRow(query, id).Scan(
		&review.ID,
		&review.Title,
		&review.EventId,
		&review.AuthorId,
		&review.Body,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrReviewNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidReviewId
	}

	return review, nil
}

// retrieves review(s) from database that match with event_id
func (r *sqlReviewRepository) GetReviewsByEventId(event_id string) ([]models.ReviewModel, error) {
	query := `SELECT
				id,
				title,
				event_id,
				author_id,
				body,
				created_at,
				updated_at

			FROM public.reviews WHERE event_id = $1`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewModel
	for rows.Next() {
		var review models.ReviewModel
		if err := rows.Scan(
			&review.ID,
			&review.Title,
			&review.EventId,
			&review.AuthorId,
			&review.Body,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// updates a review in the database
func (r *sqlReviewRepository) UpdateReview(review *models.ReviewModel) error {
	review.BeforeUpdate()

	query := `UPDATE public.reviews SET 
		title = $1,
		event_id = $2,
		author_id = $3,
		body = $4,
		updated_at = $5
		
		WHERE id = $6`

	// prevents a partial review from being submitted
	if review.CreatedAt.Unix() == 0 {
		return fmt.Errorf("unable to update a review that was not loaded from the database")
	}

	rs, err := r.database.Exec(
		query,
		review.Title,
		review.EventId,
		review.AuthorId,
		review.Body,
		review.UpdatedAt,
		review.ID,
	)
	if err != nil {
		return err
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

// deletes a review from the database.
func (r *sqlReviewRepository) DeleteReview(id string) error {
	query := `DELETE FROM public.reviews WHERE id = $1`

	rs, err := r.database.Exec(query, id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrReviewNotFound
	}

	return nil
}

// errors
var (
	ErrReviewNotFound  = errors.New("review not found")
	ErrInvalidReviewId = errors.New("invalid review id")
)
