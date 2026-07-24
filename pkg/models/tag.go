package models

// TagModel represents the tag data stored in the database.
type TagModel struct {
	ModelLifecycles `json:"-"`
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
}

func (m *TagModel) BeforeCreate() error {
	return nil
}

func (m *TagModel) AfterCreate() error {
	return nil
}

func (m *TagModel) BeforeUpdate() error {
	return nil
}

func (m *TagModel) AfterUpdate() error {
	return nil
}

func (m *TagModel) BeforeDelete() error {
	return nil
}

func (m *TagModel) AfterDelete() error {
	return nil
}
