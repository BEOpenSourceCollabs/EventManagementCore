package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventLikeRepository represents the interface for event_like-related database operations.
type EventLikeRepository interface {
	CreateEventLike(obj *models.EventLikeModel) error
	GetEventLike(event_id string, user_id string) (*models.EventLikeModel, error)
	GetUsersByEventId(event_id string) ([]models.UserModel, error)
	GetEventsByUserId(user_id string) ([]models.EventModel, error)
	DeleteEventLike(event_id string, user_id string) error
}

type sqlEventLikeRepository struct {
	database *sql.DB
}

func NewSQLEventLikeRepository(database *sql.DB) EventLikeRepository {
	return &sqlEventLikeRepository{database: database}
}

// inserts a new event_like record into the database
func (r *sqlEventLikeRepository) CreateEventLike(obj *models.EventLikeModel) error {
	obj.BeforeCreate()

	query := `INSERT INTO public.event_likes (event_id, user_id)
		VALUES ($1, $2)`

	_, err := r.database.Exec(query, obj.EventId, obj.UserId)
	if err != nil {
		return fmt.Errorf("failed to create event_like: %w", err)
	}

	obj.AfterCreate()
	return nil
}

// retrieves an event_like record from the database using event_id and user_id, if it exists
func (r *sqlEventLikeRepository) GetEventLike(event_id string, user_id string) (*models.EventLikeModel, error) {
	query := `SELECT
				event_id,
				user_id

			FROM public.event_likes WHERE event_id = $1 AND user_id = $2`

	obj := &models.EventLikeModel{}
	err := r.database.QueryRow(query, event_id, user_id).Scan(
		&obj.EventId,
		&obj.UserId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventLikeNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventLikeId
	}

	return obj, nil
}

// retrieves users from database that have liked the given event
func (r *sqlEventLikeRepository) GetUsersByEventId(event_id string) ([]models.UserModel, error) {
	query := `SELECT 
				id,
				username,
				email,
				password,
				first_name,
				last_name,
				birth_date,
				role,
				verified,
				about,
				created_at,
				updated_at,
				google_id,
				avatar_url 
	
			FROM public.users 
			WHERE id IN (
				SELECT user_id
				FROM public.event_likes 
				WHERE event_id = $1
			)`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserModel
	for rows.Next() {
		var user models.UserModel
		if err := rows.Scan(
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
			&user.GoogleId,
			&user.AvatarUrl,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// retrieves events from database that match with user_id
func (r *sqlEventLikeRepository) GetEventsByUserId(user_id string) ([]models.EventModel, error) {
	query := `SELECT
				id,
				name,
				organizer_id,
				description,
				start_date,
				end_date,
				is_paid,
				event_type,
				country,
				city, 
				slug,
				likes,
				follows,
				attendees,
				created_at,
				updated_at

			FROM public.events 
			WHERE id IN (
				SELECT event_id
				FROM public.event_likes 
				WHERE user_id = $1
			)`

	rows, err := r.database.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventModel
	for rows.Next() {
		var event models.EventModel
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.OrganizerId,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.IsPaid,
			&event.EventType,
			&event.Country,
			&event.City,
			&event.Slug,
			&event.Likes,
			&event.Follows,
			&event.Attendees,
			&event.CreatedAt,
			&event.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// deletes an event_like record from the database.
func (r *sqlEventLikeRepository) DeleteEventLike(event_id string, user_id string) error {
	query := `DELETE FROM public.event_likes WHERE event_id = $1 AND user_id = $2`

	rs, err := r.database.Exec(query, event_id, user_id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventLikeNotFound
	}

	return nil
}

// errors
var (
	ErrEventLikeNotFound  = errors.New("event_like not found")
	ErrInvalidEventLikeId = errors.New("invalid event_like id")
)
