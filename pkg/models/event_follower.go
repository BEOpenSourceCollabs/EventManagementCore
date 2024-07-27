package models

// EventFollowerModel represents the event_follower data stored in the database.
type EventFollowerModel struct {
	ModelLifecycles `json:"-"`
	EventId         string `db:"event_id" json:"event_id"`
	FollowerId      string `db:"follower_id" json:"follower_id"`
}

func (m *EventFollowerModel) BeforeCreate() error {
	return nil
}

func (m *EventFollowerModel) AfterCreate() error {
	return nil
}

func (m *EventFollowerModel) BeforeUpdate() error {
	return nil
}

func (m *EventFollowerModel) AfterUpdate() error {
	return nil
}

func (m *EventFollowerModel) BeforeDelete() error {
	return nil
}

func (m *EventFollowerModel) AfterDelete() error {
	return nil
}
