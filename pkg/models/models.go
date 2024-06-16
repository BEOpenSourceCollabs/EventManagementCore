package models

import (
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/logger"
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
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle BeforeCreate()")
	return nil
}

func (m *Model) AfterCreate() error {
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle AfterCreate()")
	return nil
}

func (m *Model) BeforeUpdate() error {
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle BeforeUpdate()")
	return nil
}

func (m *Model) AfterUpdate() error {
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle AfterUpdate()")
	return nil
}

func (m *Model) BeforeDelete() error {
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle BeforeDelete()")
	return nil
}

func (m *Model) AfterDelete() error {
	logger.AppLogger.InfoF("Model lifecycle", "unoverridden lifecycle AfterDelete()")
	return nil
}
