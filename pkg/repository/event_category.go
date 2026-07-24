package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
)

// EventCategoryRepository represents the interface for event_category-related database operations.
type EventCategoryRepository interface {
	CreateEventCategory(obj *models.EventCategoryModel) error
	GetEventCategory(event_id string, category_id string) (*models.EventCategoryModel, error)
	GetCategoriesByEventId(event_id string) ([]models.CategoryModel, error)
	GetEventsByCategoryId(category_id string) ([]models.EventModel, error)
	DeleteEventCategory(event_id string, category_id string) error
}

type sqlEventCategoryRepository struct {
	database *sql.DB
}

func NewSQLEventCategoryRepository(database *sql.DB) EventCategoryRepository {
	return &sqlEventCategoryRepository{database: database}
}

// inserts a new event_category into the database
func (r *sqlEventCategoryRepository) CreateEventCategory(obj *models.EventCategoryModel) error {
	obj.BeforeCreate()

	query := `INSERT INTO public.event_categories (event_id, category_id)
		VALUES ($1, $2)`

	_, err := r.database.Exec(query, obj.EventId, obj.CategoryId)
	if err != nil {
		return fmt.Errorf("failed to create event_category: %w", err)
	}

	obj.AfterCreate()
	return nil
}

// retrieves an event_category from the database using its event_id and category_id, if it exists
func (r *sqlEventCategoryRepository) GetEventCategory(event_id string, category_id string) (*models.EventCategoryModel, error) {
	query := `SELECT
				event_id,
				category_id

			FROM public.event_categories WHERE event_id = $1 AND category_id = $2`

	obj := &models.EventCategoryModel{}
	err := r.database.QueryRow(query, event_id, category_id).Scan(
		&obj.EventId,
		&obj.CategoryId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEventCategoryNotFound
		}
		if reflect.TypeOf(err) == reflect.TypeOf(&net.OpError{}) {
			return nil, ErrRepoConnErr
		}
		return nil, ErrInvalidEventCategoryId
	}

	return obj, nil
}

// retrieves categories from the database that are associated with the given event
func (r *sqlEventCategoryRepository) GetCategoriesByEventId(event_id string) ([]models.CategoryModel, error) {
	query := `SELECT 
				id,
				name 
	
			FROM public.categories 
			WHERE id IN (
				SELECT category_id
				FROM public.event_categories 
				WHERE event_id = $1
			)`

	rows, err := r.database.Query(query, event_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.CategoryModel
	for rows.Next() {
		var category models.CategoryModel
		if err := rows.Scan(
			&category.ID,
			&category.Name,
		); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// retrieves events from database that match with category_id
func (r *sqlEventCategoryRepository) GetEventsByCategoryId(category_id string) ([]models.EventModel, error) {
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
				FROM public.event_categories 
				WHERE category_id = $1
			)`

	rows, err := r.database.Query(query, category_id)
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

// deletes an event_category from the database.
func (r *sqlEventCategoryRepository) DeleteEventCategory(event_id string, category_id string) error {
	query := `DELETE FROM public.event_categories WHERE event_id = $1 AND category_id = $2`

	rs, err := r.database.Exec(query, event_id, category_id)
	if err != nil {
		return err
	}

	if affected, err := rs.RowsAffected(); affected < 1 {
		if err != nil {
			return err
		}
		return ErrEventCategoryNotFound
	}

	return nil
}

// errors
var (
	ErrEventCategoryNotFound  = errors.New("event_category not found")
	ErrInvalidEventCategoryId = errors.New("invalid event_category id")
)
