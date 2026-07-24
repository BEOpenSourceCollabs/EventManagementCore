package models

import "time"

// EventAttendeeModel represents the event_attendee data stored in the database.
type EventAttendeeModel struct {
	ModelLifecycles `json:"-"`
	EventId         string    `db:"event_id" json:"event_id"`
	AttendeeId      string    `db:"attendee_id" json:"attendee_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func (m *EventAttendeeModel) BeforeCreate() error {
	return nil
}

func (m *EventAttendeeModel) AfterCreate() error {
	return nil
}

func (m *EventAttendeeModel) BeforeUpdate() error {
	return nil
}

func (m *EventAttendeeModel) AfterUpdate() error {
	return nil
}

func (m *EventAttendeeModel) BeforeDelete() error {
	return nil
}

func (m *EventAttendeeModel) AfterDelete() error {
	return nil
}
