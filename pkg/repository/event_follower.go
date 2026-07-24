package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventFollowerRepository represents the interface for event_follower-related database operations.
type EventFollowerRepository interface {
	CreateEventFollower(obj *models.EventFollowerModel) error
	GetEventFollower(event_id string, follower_id string) (*models.EventFollowerModel, error)
	GetFollowersByEventId(event_id string) ([]models.UserModel, error)
	GetEventsByFollowerId(follower_id string) ([]models.EventModel, error)
	DeleteEventFollower(event_id string, follower_id string) error
}

type sqlEventFollowerRepository struct {
	database *sql.DB
}

func NewSQLEventFollowerRepository(database *sql.DB) EventFollowerRepository {
	return &sqlEventFollowerRepository{database: database}
}

// inserts a new event_follower record into the database
func (r *sqlEventFollowerRepository) CreateEventFollower(obj *models.EventFollowerModel) error {
	obj.BeforeCreate()

	query := `INSERT INTO public.event_followers (event_id, follower_id)
		VALUES ($1, $2)`

	_, err := r.database.Exec(query, obj.EventId, obj.FollowerId)
	if err != nil {
		return fmt.Errorf("failed to create event_follower: %w", err)
	}

	obj.AfterCreate()
	return nil
}

// retrieves an event_follower record from the database using its event_id and follower_id, if it exists
func (r *sqlEventFollowerRepository) GetEventFollower(event_id string, follower_id string) (*models.EventFollowerModel, error) {
	query := `SELECT
				event_id,
				follower_id

			FROM public.event_followers WHERE event_id = $1 AND follower_id = $2`

	obj := &models.EventFollowerModel{}
	err := r.database.QueryRow(query, event_id, follower_id).Scan(
		&obj.EventId,
		&obj.FollowerId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventFollowerNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventFollowerId
	}

	return obj, nil
}

// retrieves users from database that are followers of the given event
func (r *sqlEventFollowerRepository) GetFollowersByEventId(event_id string) ([]models.UserModel, error) {
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
				SELECT follower_id
				FROM public.event_followers 
				WHERE event_id = $1
			)`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.UserModel
	for rows.Next() {
		var follower models.UserModel
		if err := rows.Scan(
			&follower.ID,
			&follower.Username,
			&follower.Email,
			&follower.Password,
			&follower.FirstName,
			&follower.LastName,
			&follower.BirthDate,
			&follower.Role,
			&follower.Verified,
			&follower.About,
			&follower.CreatedAt,
			&follower.UpdatedAt,
			&follower.GoogleId,
			&follower.AvatarUrl,
		); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

// retrieves events from database that match with follower_id
func (r *sqlEventFollowerRepository) GetEventsByFollowerId(follower_id string) ([]models.EventModel, error) {
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
				FROM public.event_followers 
				WHERE follower_id = $1
			)`

	rows, err := r.database.Query(query, follower_id)
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

// deletes an event_follower record from the database.
func (r *sqlEventFollowerRepository) DeleteEventFollower(event_id string, follower_id string) error {
	query := `DELETE FROM public.event_followers WHERE event_id = $1 AND follower_id = $2`

	rs, err := r.database.Exec(query, event_id, follower_id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventFollowerNotFound
	}

	return nil
}

// errors
var (
	ErrEventFollowerNotFound  = errors.New("event_follower not found")
	ErrInvalidEventFollowerId = errors.New("invalid event_follower id")
)
