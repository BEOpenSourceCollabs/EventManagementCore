package models

import (
	"database/sql"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

// EventModel represents the event data stored in the database.
type EventModel struct {
	Model
	Name        string          `db:"name" json:"name"`
	OrganizerId string          `db:"organizer_id" json:"organizer_id"`
	Description sql.NullString  `db:"description" json:"description"`
	StartDate   time.Time       `db:"start_date" json:"start_date"`
	EndDate     time.Time       `db:"end_date" json:"end_date"`
	IsPaid      bool            `db:"is_paid" json:"is_paid"`
	EventType   types.EventType `db:"event_type" json:"event_type"`
	Country     sql.NullString  `db:"country" json:"country"`
	City        sql.NullString  `db:"city" json:"city"`
	Slug        string          `db:"slug" json:"slug"`
	Likes       int64           `db:"likes" json:"likes"`
	Follows     int64           `db:"follows" json:"follows"`
	Attendees   int64           `db:"attendees" json:"attendees"`
}

// BeforeCreate overrides model lifecycle hook
func (m *EventModel) BeforeCreate() error {
	m.Likes = 0
	m.Follows = 0
	m.Attendees = 0
	// TODO: generate value for slug
	return nil
}

// BeforeUpdate overrides model lifecycle hook, updating the updated_at time.
func (m *EventModel) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
