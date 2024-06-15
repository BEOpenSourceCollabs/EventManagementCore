package models

import (
	"database/sql"
	"time"
)

// EventModel represents the event data stored in the database.
type EventModel struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	OrganizerID string         `db:"organizer_id"`
	Description sql.NullString `db:"description"`
	StartDate   time.Time      `db:"start_date"`
	EndDate     time.Time      `db:"end_date"`
	IsPaid      bool           `db:"is_paid"`
	EventType   string         `db:"event_type"`
	Country     sql.NullString `db:"country"`
	City        sql.NullString `db:"city"`
	Slug        string         `db:"slug"`
	Likes       int64          `db:"likes"`
	Follows     int64          `db:"follows"`
	Attendees   int64          `db:"attendees"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}
