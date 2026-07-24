package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventTagRepository represents the interface for event_tag-related database operations.
type EventTagRepository interface {
	CreateEventTag(obj *models.EventTagModel) error
	GetEventTag(event_id string, tag_id string) (*models.EventTagModel, error)
	GetTagsByEventId(event_id string) ([]models.TagModel, error)
	GetEventsByTagId(tag_id string) ([]models.EventModel, error)
	DeleteEventTag(event_id string, tag_id string) error
}

type sqlEventTagRepository struct {
	database *sql.DB
}

func NewSQLEventTagRepository(database *sql.DB) EventTagRepository {
	return &sqlEventTagRepository{database: database}
}

// inserts a new event_tag into the database
func (r *sqlEventTagRepository) CreateEventTag(obj *models.EventTagModel) error {
	obj.BeforeCreate()

	query := `INSERT INTO public.event_tags (event_id, tag_id)
		VALUES ($1, $2)`

	_, err := r.database.Exec(query, obj.EventId, obj.TagId)
	if err != nil {
		return fmt.Errorf("failed to create event_tag: %w", err)
	}

	obj.AfterCreate()
	return nil
}

// retrieves an event_tag from the database using its event_id and tag_id, if it exists
func (r *sqlEventTagRepository) GetEventTag(event_id string, tag_id string) (*models.EventTagModel, error) {
	query := `SELECT
				event_id,
				tag_id

			FROM public.event_tags WHERE event_id = $1 AND tag_id = $2`

	obj := &models.EventTagModel{}
	err := r.database.QueryRow(query, event_id, tag_id).Scan(
		&obj.EventId,
		&obj.TagId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventTagNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventTagId
	}

	return obj, nil
}

// retrieves tags from the database that are associated with the given event
func (r *sqlEventTagRepository) GetTagsByEventId(event_id string) ([]models.TagModel, error) {
	query := `SELECT 
				id,
				name 
	
			FROM public.tags 
			WHERE id IN (
				SELECT tag_id
				FROM public.event_tags 
				WHERE event_id = $1
			)`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []models.TagModel
	for rows.Next() {
		var tag models.TagModel
		if err := rows.Scan(
			&tag.ID,
			&tag.Name,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// retrieves events from database that match with tag_id
func (r *sqlEventTagRepository) GetEventsByTagId(tag_id string) ([]models.EventModel, error) {
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
				FROM public.event_tags 
				WHERE tag_id = $1
			)`

	rows, err := r.database.Query(query, tag_id)
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

// deletes an event_tag from the database.
func (r *sqlEventTagRepository) DeleteEventTag(event_id string, tag_id string) error {
	query := `DELETE FROM public.event_tags WHERE event_id = $1 AND tag_id = $2`

	rs, err := r.database.Exec(query, event_id, tag_id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventTagNotFound
	}

	return nil
}

// errors
var (
	ErrEventTagNotFound  = errors.New("event_tag not found")
	ErrInvalidEventTagId = errors.New("invalid event_tag id")
)
