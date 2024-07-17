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
	reviews = []models.ReviewModel{
		{
			Title: "Great Event",
			Body:  "This was an amazing event!",
		},
		{
			Title: "Good Experience",
			Body:  "I enjoyed this event, but there's room for improvement.",
		},
		{
			Title: "Disappointing",
			Body:  "The event didn't meet my expectations.",
		},
	}
)

func TestReviewRepository_KitchenSink(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	reviewRepo := repository.NewSQLReviewRepository(db)
	userRepo := repository.NewSQLUserRepository(db)
	eventRepo := repository.NewSQLEventRepository(db)

	author := models.UserModel{
		Username: "author0",
		FirstName: sql.NullString{
			String: "author",
			Valid:  true,
		},
		LastName: sql.NullString{
			String: "0",
			Valid:  true,
		},
		Email:    "author0@example.co.uk",
		Verified: true,
		Role:     types.UserRole,
	}

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

	event := models.EventModel{
		Name:        "Coding Meetup",
		Type:        types.BothEventType,
		Description: "A test networking event",
		StartDate:   time.Now().Add(time.Hour + 48),
		EndDate:     time.Now().Add(time.Hour + 49),
		IsPaid:      false,
		CountryISO:  "GB",
		City:        "London",
	}

	t.Run("Create review author", func(t *testing.T) {
		if err := userRepo.CreateUser(&author); err != nil {
			t.Errorf("expected no error when creating review author but got %v", err)
		}
	})

	t.Run("Create event organizer", func(t *testing.T) {
		if err := userRepo.CreateUser(&organizer); err != nil {
			t.Errorf("expected no error when creating organizer but got %v", err)
		}
	})

	t.Run("Create event", func(t *testing.T) {
		event.Organizer = organizer.ID
		if err := eventRepo.CreateEvent(&event); err != nil {
			t.Errorf("expected no error when creating event but got %v", err)
		}
	})

	t.Run("Create reviews", func(t *testing.T) {
		for i := range reviews {
			reviews[i].AuthorID = author.ID // FKC requires a valid author id
			reviews[i].EventID = event.ID   // FKC requires a valid event id

			if err := reviewRepo.CreateReview(&reviews[i]); err != nil {
				t.Errorf("expected no error when creating review but got %v", err)
			}
			if reviews[i].ID == "" {
				t.Error("expected review ID to be set after creation")
			}
		}
	})

	t.Run("Update review without loading first", func(t *testing.T) {
		if err := reviewRepo.UpdateReview(&models.ReviewModel{Model: models.Model{ID: reviews[0].ID}, Title: "Updated"}); err == nil {
			t.Error("expected error while attempting to update a review without fully loading the review first")
		}
	})

	t.Run("Update review", func(t *testing.T) {
		reviews[0].Title = "Updated Title"
		reviews[0].Body = "Updated body content"
		if err := reviewRepo.UpdateReview(&reviews[0]); err != nil {
			t.Errorf("expected no error when updating review but got %v", err)
		}
		loaded, err := reviewRepo.GetReviewByID(reviews[0].ID)
		if err != nil {
			t.Errorf("expected no error when getting review by id but got %v", err)
		}
		if loaded.Title != "Updated Title" {
			t.Errorf("expected loaded review's title to be updated to 'Updated Title' but was '%s'", loaded.Title)
		}
		if loaded.Body != "Updated body content" {
			t.Errorf("expected loaded review's body to be updated but was '%s'", loaded.Body)
		}
	})

	t.Run("Get review by ID", func(t *testing.T) {
		loaded, err := reviewRepo.GetReviewByID(reviews[1].ID)
		if err != nil {
			t.Errorf("expected no error when getting review by ID but got %v", err)
		}
		if loaded.ID != reviews[1].ID {
			t.Errorf("expected loaded review's id to match '%s' but was '%s'", reviews[1].ID, loaded.ID)
		}
	})

	t.Run("Delete reviews", func(t *testing.T) {
		for _, review := range reviews {
			if err := reviewRepo.DeleteReview(review.ID); err != nil {
				t.Errorf("expected no error when deleting review but got %v", err)
			}
			deletedReview, err := reviewRepo.GetReviewByID(review.ID)
			if deletedReview != nil {
				t.Errorf("expected deleted review to be nil but was %v", deletedReview)
			}
			if err == nil {
				t.Errorf("expected error when attempting to get review by id after deletion")
			}
		}
	})
}
