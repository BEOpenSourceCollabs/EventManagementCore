package models

import (
	"time"
)

/*
CREATE TABLE IF NOT EXISTS public.events (
   name VARCHAR(500) NOT NULL,
   organizer_id UUID NOT NULL,
   description VARCHAR(1000),
   start_date TIMESTAMPTZ NOT NULL,
   end_date TIMESTAMPTZ NOT NULL,
   is_paid BOOLEAN default 'false',
   event_type event_type NOT NULL DEFAULT 'online',
   country varchar(5),
   city varchar(50),
   slug text not null,
   likes INT NOT NULL DEFAULT 0,
   follows INT NOT NULL DEFAULT 0,
   attendees INT NOT NULL DEFAULT 0,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (organizer_id) REFERENCES public.users(id) ON DELETE CASCADE
);
*/

// EventModel represents the event data stored in the database.
type EventModel struct {
	Model
	Name        string    `db:"name" json:"name"`
	Organizer   string    `db:"organizer_id" json:"organizer_id"`
	Description string    `db:"description" json:"description"`
	StartDate   time.Time `db:"start_date" json:"start_date"`
	EndDate     time.Time `db:"end_date" json:"end_date"`
	IsPaid      string    `db:"is_paid" json:"is_paid"`
	Type        string    `db:"event_type" json:"event_type"`
	Country     string    `db:"country" json:"country"`
	City        string    `db:"city" json:"city"`
	Slug        string    `db:"slug" json:"slug"`
	Likes       int64     `db:"likes" json:"likes"`
	Follows     int64     `db:"follows" json:"follows"`
	Attendees   int64     `db:"attendees" json:"attendees"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *EventModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
