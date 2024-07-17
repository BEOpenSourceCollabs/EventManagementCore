package models

import "time"

type ReviewModel struct {
	Model
	Title    string `db:"title" json:"title"`
	EventID  string `db:"event_id" json:"event_id"`
	AuthorID string `db:"author_id" json:"author_id"`
	Body     string `db:"body" json:"body"`
}

// BeforeUpdated overrides model lifecycle hook, updating the updated_at time.
func (m *ReviewModel) BeforeUpdated() error {
	m.UpdatedAt = time.Now()
	return nil
}
