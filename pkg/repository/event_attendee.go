package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventAttendeeRepository represents the interface for event_attendee-related database operations.
type EventAttendeeRepository interface {
	CreateEventAttendee(obj *models.EventAttendeeModel) error
	GetEventAttendee(event_id string, attendee_id string) (*models.EventAttendeeModel, error)
	GetAttendeesByEventId(event_id string) ([]models.UserModel, error)
	GetEventsByAttendeeId(attendee_id string) ([]models.EventModel, error)
	DeleteEventAttendee(event_id string, attendee_id string) error
}

type sqlEventAttendeeRepository struct {
	database *sql.DB
}

func NewSQLEventAttendeeRepository(database *sql.DB) EventAttendeeRepository {
	return &sqlEventAttendeeRepository{database: database}
}

// inserts a new event_attendee record into the database
func (r *sqlEventAttendeeRepository) CreateEventAttendee(obj *models.EventAttendeeModel) error {
	obj.BeforeCreate()

	query := `INSERT INTO public.event_attendees (event_id, attendee_id)
		VALUES ($1, $2)`

	_, err := r.database.Exec(query, obj.EventId, obj.AttendeeId)
	if err != nil {
		return fmt.Errorf("failed to create event_attendee: %w", err)
	}

	obj.AfterCreate()
	return nil
}

// retrieves an event_attendee record from the database using event_id and attendee_id, if it exists
func (r *sqlEventAttendeeRepository) GetEventAttendee(event_id string, attendee_id string) (*models.EventAttendeeModel, error) {
	query := `SELECT
				event_id,
				attendee_id,
				created_at

			FROM public.event_attendees WHERE event_id = $1 AND attendee_id = $2`

	obj := &models.EventAttendeeModel{}
	err := r.database.QueryRow(query, event_id, attendee_id).Scan(
		&obj.EventId,
		&obj.AttendeeId,
		&obj.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventAttendeeNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventAttendeeId
	}

	return obj, nil
}

// retrieves users from database that are attendees of the given event
func (r *sqlEventAttendeeRepository) GetAttendeesByEventId(event_id string) ([]models.UserModel, error) {
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
				SELECT attendee_id
				FROM public.event_attendees 
				WHERE event_id = $1
			)`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []models.UserModel
	for rows.Next() {
		var attendee models.UserModel
		if err := rows.Scan(
			&attendee.ID,
			&attendee.Username,
			&attendee.Email,
			&attendee.Password,
			&attendee.FirstName,
			&attendee.LastName,
			&attendee.BirthDate,
			&attendee.Role,
			&attendee.Verified,
			&attendee.About,
			&attendee.CreatedAt,
			&attendee.UpdatedAt,
			&attendee.GoogleId,
			&attendee.AvatarUrl,
		); err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

// retrieves events from database that match with attendee_id
func (r *sqlEventAttendeeRepository) GetEventsByAttendeeId(attendee_id string) ([]models.EventModel, error) {
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
				FROM public.event_attendees 
				WHERE attendee_id = $1
			)`

	rows, err := r.database.Query(query, attendee_id)
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

// deletes an event_attendee record from the database.
func (r *sqlEventAttendeeRepository) DeleteEventAttendee(event_id string, attendee_id string) error {
	query := `DELETE FROM public.event_attendees WHERE event_id = $1 AND attendee_id = $2`

	rs, err := r.database.Exec(query, event_id, attendee_id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventAttendeeNotFound
	}

	return nil
}

// errors
var (
	ErrEventAttendeeNotFound  = errors.New("event_attendee not found")
	ErrInvalidEventAttendeeId = errors.New("invalid event_attendee id")
)
