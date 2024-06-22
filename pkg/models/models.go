package models

import (
	"time"
)

// ModelLifecycles should be invoked from the repository before and after create, update and deletion of a model which implements the interface.
type ModelLifecycles interface {
	BeforeCreate() error
	AfterCreate() error
	BeforeUpdate() error
	AfterUpdate() error
	BeforeDelete() error
	AfterDelete() error
}

// Model provides default fields and lifecycle functions to models.
type Model struct {
	ModelLifecycles `json:"-"`
	ID              string    `db:"id" json:"id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

func (m *Model) BeforeCreate() error {
	return nil
}

func (m *Model) AfterCreate() error {
	return nil
}

func (m *Model) BeforeUpdate() error {
	return nil
}

func (m *Model) AfterUpdate() error {
	return nil
}

func (m *Model) BeforeDelete() error {
	return nil
}

func (m *Model) AfterDelete() error {
	return nil
}
