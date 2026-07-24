package models

import "time"

// ReviewModel represents the review data stored in the database.
type ReviewModel struct {
	Model
	Title    string `db:"title" json:"title"`
	EventId  string `db:"event_id" json:"event_id"`
	AuthorId string `db:"author_id" json:"author_id"`
	Body     string `db:"body" json:"body"`
}

// BeforeUpdate overrides model lifecycle hook, updating the updated_at time.
func (m *ReviewModel) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}
