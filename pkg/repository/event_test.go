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
			Type:        types.OfflineEventType,
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
		for i := range events {
			events[i].Organizer = organizer.ID // need a valid organizer for the event

			if err := eventRepo.CreateEvent(&events[i]); err != nil {
				t.Errorf("expected no error when creating event but got %v", err)
			}
		}
	})

	t.Run("Update event", func(t *testing.T) {
		events[0].Name = "Virtual Music Festival"
		events[0].Type = types.BothEventType
		if err := eventRepo.UpdateEvent(&events[0]); err != nil {
			t.Errorf("expected no error when updating event but got %v", err)
		}

		// verify updated
		loadedEvent, err := eventRepo.GetEventByID(events[0].ID)
		if err != nil {
			t.Errorf("expected no error when getting event by id '%s' but got %v", events[0].ID, err)
		}
		if loadedEvent.Name != events[0].Name {
			t.Errorf("expected '%s' as updated event name but was '%s'", events[0].Name, loadedEvent.Name)
		}
		if loadedEvent.Type != events[0].Type {
			t.Errorf("expected '%s' as updated event type but was '%s'", events[0].Type, loadedEvent.Type)
		}
	})

	t.Run("Delete events", func(t *testing.T) {
		for _, event := range events {
			if err := eventRepo.DeleteEvent(event.ID); err != nil {
				t.Errorf("expected no error when deleting event but got %v", err)
			}
			devnt, err := eventRepo.GetEventByID(event.ID)
			if devnt != nil {
				t.Errorf("expected deleted event to be nil but was %v", devnt)
			}
			if err == nil {
				t.Errorf("expected error when attempting to get event by id after deletion")
			}
		}
	})
}
