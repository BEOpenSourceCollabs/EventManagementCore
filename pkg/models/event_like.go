package models

// EventLikeModel represents the event_like data stored in the database.
type EventLikeModel struct {
	ModelLifecycles `json:"-"`
	EventId         string `db:"event_id" json:"event_id"`
	UserId          string `db:"user_id" json:"user_id"`
}

func (m *EventLikeModel) BeforeCreate() error {
	return nil
}

func (m *EventLikeModel) AfterCreate() error {
	return nil
}

func (m *EventLikeModel) BeforeUpdate() error {
	return nil
}

func (m *EventLikeModel) AfterUpdate() error {
	return nil
}

func (m *EventLikeModel) BeforeDelete() error {
	return nil
}

func (m *EventLikeModel) AfterDelete() error {
	return nil
}
