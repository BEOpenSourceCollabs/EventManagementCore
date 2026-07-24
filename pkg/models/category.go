package models

// CategoryModel represents the category data stored in the database.
type CategoryModel struct {
	ModelLifecycles `json:"-"`
	ID              string `db:"id" json:"id"`
	Name            string `db:"name" json:"name"`
}

func (m *CategoryModel) BeforeCreate() error {
	return nil
}

func (m *CategoryModel) AfterCreate() error {
	return nil
}

func (m *CategoryModel) BeforeUpdate() error {
	return nil
}

func (m *CategoryModel) AfterUpdate() error {
	return nil
}

func (m *CategoryModel) BeforeDelete() error {
	return nil
}

func (m *CategoryModel) AfterDelete() error {
	return nil
}
