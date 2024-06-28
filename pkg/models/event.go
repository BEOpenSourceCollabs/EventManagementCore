package models

import (
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

// EventModel represents the event data stored in the database.
type EventModel struct {
	Model
	Name        string          `db:"name" json:"name"`
	Organizer   string          `db:"organizer_id" json:"organizer_id"`
	Description string          `db:"description" json:"description"`
	StartDate   time.Time       `db:"start_date" json:"start_date"`
	EndDate     time.Time       `db:"end_date" json:"end_date"`
	IsPaid      bool            `db:"is_paid" json:"is_paid"`
	Type        types.EventType `db:"event_type" json:"event_type"`
	CountryISO  string          `db:"country" json:"country"`
	City        string          `db:"city" json:"city"`
	Slug        string          `db:"slug" json:"slug"`
	Likes       int64           `db:"likes" json:"likes"`
	Follows     int64           `db:"follows" json:"follows"`
	Attendees   int64           `db:"attendees" json:"attendees"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *EventModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
