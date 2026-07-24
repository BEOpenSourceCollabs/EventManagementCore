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
	GetEventsByOrganizerId(org_id string) ([]models.EventModel, error)
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

// inserts a new event into the database
func (r *sqlEventRepository) CreateEvent(event *models.EventModel) error {
	event.BeforeCreate()

	query := `INSERT INTO public.events (name, organizer_id, description, start_date, end_date, is_paid, event_type, country, city, slug, likes, follows, attendees)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`

	err := r.database.QueryRow(query, event.Name, event.OrganizerId, event.Description, event.StartDate, event.EndDate, event.IsPaid, event.EventType, event.Country, event.City, event.Slug, event.Likes, event.Follows, event.Attendees).Scan(&event.ID)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	event.AfterCreate()
	return nil
}

// retrieves an event from the database by its unique ID, if it exists
func (r *sqlEventRepository) GetEventByID(id string) (*models.EventModel, error) {
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

			FROM public.events WHERE id = $1`

	event := &models.EventModel{}
	err := r.database.QueryRow(query, id).Scan(
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

// retrieves event(s) from database that match with organizer_id
func (r *sqlEventRepository) GetEventsByOrganizerId(org_id string) ([]models.EventModel, error) {
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

			FROM public.events WHERE organizer_id = $1`

	rows, err := r.database.Query(query, org_id)
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

// updates an event in database
func (r *sqlEventRepository) UpdateEvent(event *models.EventModel) error {
	event.BeforeUpdate()

	query := `UPDATE public.events SET 
		name = $1, 
		organizer_id = $2, 
		description = $3, 
		start_date = $4, 
		end_date = $5, 
		is_paid = $6, 
		event_type = $7, 
		country = $8, 
		city = $9, 
		slug = $10, 
		likes = $11, 
		follows = $12, 
		attendees = $13, 
		updated_at = $14 
		
		WHERE id = $15`

	// prevents a partial event from being submitted
	if event.CreatedAt.Unix() == 0 {
		return fmt.Errorf("unable to update an event that was not loaded from the database")
	}

	rs, err := r.database.Exec(
		query,
		event.Name,
		event.OrganizerId,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.IsPaid,
		event.EventType,
		event.Country,
		event.City,
		event.Slug,
		event.Likes,
		event.Follows,
		event.Attendees,
		event.UpdatedAt,
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

// deletes an event from the database.
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

// errors
var (
	ErrEventNotFound  = errors.New("event not found")
	ErrInvalidEventId = errors.New("invalid event id")
)
