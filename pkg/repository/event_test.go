package repository_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/test"
)

var (
	events = []models.EventModel{
		{
			Name:        "Music Festival",
			Type:        "Music",
			Organizer:   "",
			Description: "A test music festival",
			StartDate:   time.Now().Add(time.Hour + 24),
			EndDate:     time.Now().Add(time.Hour + 25),
			IsPaid:      true,
			Country:     "England",
			City:        "London",
		},
		{
			Name:        "Coding Meetup",
			Type:        "Networking",
			Organizer:   "",
			Description: "A test networking event",
			StartDate:   time.Now().Add(time.Hour + 48),
			EndDate:     time.Now().Add(time.Hour + 49),
			IsPaid:      false,
			Country:     "England",
			City:        "London",
		},
	}
)

func TestEventRepository_KitchenSink(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(persist.DatabaseConfiguration{
		User:     "postgres",
		Password: "postgres",
		Database: "event-mgmt-db",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Clean up the container
	defer func() {
		if err := container.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()
	eventRepo := repository.NewSQLEventRepository(db)

	t.Run("Create events", func(t *testing.T) {
		for _, event := range events {
			if err := eventRepo.CreateEvent(&event); err != nil {
				t.Errorf("expected no error when creating event but got %v", err)
			}
		}
	})
}
