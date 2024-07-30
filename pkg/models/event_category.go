package models

// EventCategoryModel represents the event_category data stored in the database.
type EventCategoryModel struct {
	ModelLifecycles `json:"-"`
	EventId         string `db:"event_id" json:"event_id"`
	CategoryId      string `db:"category_id" json:"category_id"`
}

func (m *EventCategoryModel) BeforeCreate() error {
	return nil
}

func (m *EventCategoryModel) AfterCreate() error {
	return nil
}

func (m *EventCategoryModel) BeforeUpdate() error {
	return nil
}

func (m *EventCategoryModel) AfterUpdate() error {
	return nil
}

func (m *EventCategoryModel) BeforeDelete() error {
	return nil
}

func (m *EventCategoryModel) AfterDelete() error {
	return nil
}
