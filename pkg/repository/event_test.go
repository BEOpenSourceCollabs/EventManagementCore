package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/test"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

var (
	events = []models.EventModel{
		{
			Name:        "Music Festival",
			Type:        types.OnlineEventType,
			Description: "A test music festival",
			StartDate:   time.Now().Add(time.Hour + 24),
			EndDate:     time.Now().Add(time.Hour + 25),
			IsPaid:      true,
			CountryISO:  "GB",
			City:        "London",
		},
		{
			Name:        "Coding Meetup",
			Type:        types.BothEventType,
			Description: "A test networking event",
			StartDate:   time.Now().Add(time.Hour + 48),
			EndDate:     time.Now().Add(time.Hour + 49),
			IsPaid:      false,
			CountryISO:  "GB",
			City:        "London",
		},
	}
)

func TestEventRepository_KitchenSink(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	eventRepo := repository.NewSQLEventRepository(db)
	userRepo := repository.NewSQLUserRepository(db)

	organizer := models.UserModel{
		Username: "org0",
		FirstName: sql.NullString{
			String: "org",
			Valid:  true,
		},
		LastName: sql.NullString{
			String: "0",
			Valid:  true,
		},
		Email:    "org0@example.co.uk",
		Verified: true,
		Role:     types.OrganizerRole,
	}

	t.Run("Create event organizer", func(t *testing.T) {
		if err := userRepo.CreateUser(&organizer); err != nil {
			t.Errorf("expected no error when creating organizer but got %v", err)
		}
	})

	t.Run("Create events", func(t *testing.T) {
		for _, event := range events {
			event.Organizer = organizer.ID // need a valid organizer for the event

			if err := eventRepo.CreateEvent(&event); err != nil {
				t.Errorf("expected no error when creating event but got %v", err)
			}
		}
	})
}
