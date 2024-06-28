package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventRepository represents the interface for event-related database operations.
type EventRepository interface {
	CreateEvent(event *models.EventModel) error
	GetEventByID(id string) (*models.EventModel, error)
	UpdateEvent(event *models.EventModel) error
	DeleteEvent(id string) error
}

type sqlEventRepository struct {
	database *sql.DB
}

// NewSQLEventRepository creates and returns a new sql flavoured EventRepository instance.
func NewSQLEventRepository(database *sql.DB) EventRepository {
	return &sqlEventRepository{database: database}
}

// CreateEvent inserts a new event into the database.
func (r *sqlEventRepository) CreateEvent(event *models.EventModel) error {
	event.BeforeCreate()

	query := `INSERT INTO public.events (name, organizer_id, description, start_date, end_date, is_paid, event_type, country, city, slug, likes, follows, attendees)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	err := r.database.QueryRow(query,
		event.Name,
		event.Organizer,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.IsPaid,
		event.Type,
		event.CountryISO,
		event.City,
		event.Slug,
		event.Likes,
		event.Follows,
		event.Attendees,
	).Scan(&event.ID)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	event.AfterCreate()

	return nil
}

// GetEventByID retrieves a event from the database by its unique ID.
func (r *sqlEventRepository) GetEventByID(id string) (*models.EventModel, error) {
	query := `SELECT id, name, organizer_id, description, start_date, end_date, is_paid, event_type, country, city, slug, likes, follows, attendees, created_at, updated_at,
				FROM public.events WHERE id = $1`

	event := &models.EventModel{}
	err := r.database.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Organizer,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.IsPaid,
		&event.Type,
		&event.CountryISO,
		&event.City,
		&event.Slug,
		&event.Likes,
		&event.Follows,
		&event.Attendees,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventId
	}

	return event, nil
}

// UpdateEvent update a event in the database.
func (r *sqlEventRepository) UpdateEvent(event *models.EventModel) error {
	event.BeforeUpdate()
	query := `UPDATE public.events SET name = $1, organizer_id = $2, description = $3, start_date = $4, end_date = $5, is_paid = $6, event_type = $7, country = $8, city = $9, slug = $10, likes = $11, follows = $12, attendees = $13, WHERE id = $14`
	// This is a guard to prevent any partial event from being submitted.
	// Otherwise it would be possible to accidently empty out columns by passing empty/uninitialized values.
	if event.CreatedAt.Unix() == 0 {
		return fmt.Errorf("unable to update a event that was not loaded from the database")
	}

	rs, err := r.database.Exec(
		query,
		event.Name,
		event.Organizer,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.IsPaid,
		event.Type,
		event.CountryISO,
		event.City,
		event.Slug,
		event.Likes,
		event.Follows,
		event.Attendees,
		event.ID,
	)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventNotFound
	}

	event.AfterUpdate()

	return nil
}

// DeleteEvent delete a event from the database.
func (r *sqlEventRepository) DeleteEvent(id string) error {
	query := `DELETE FROM public.events WHERE id = $1`

	rs, err := r.database.Exec(query, id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventNotFound
	}

	return nil
}

var (
	ErrEventNotFound  = errors.New("event not found")  // ErrEventNotFound is returned when a event is not found in the database.
	ErrInvalidEventId = errors.New("invalid event id") // ErrEventNotFound is returned when a event id is invalid or malformed.
)
