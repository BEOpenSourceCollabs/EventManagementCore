package models

// EventTagModel represents the event_tag data stored in the database.
type EventTagModel struct {
	ModelLifecycles `json:"-"`
	EventId         string `db:"event_id" json:"event_id"`
	TagId           string `db:"tag_id" json:"tag_id"`
}

func (m *EventTagModel) BeforeCreate() error {
	return nil
}

func (m *EventTagModel) AfterCreate() error {
	return nil
}

func (m *EventTagModel) BeforeUpdate() error {
	return nil
}

func (m *EventTagModel) AfterUpdate() error {
	return nil
}

func (m *EventTagModel) BeforeDelete() error {
	return nil
}

func (m *EventTagModel) AfterDelete() error {
	return nil
}
