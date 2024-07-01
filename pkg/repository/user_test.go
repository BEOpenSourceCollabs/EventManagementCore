package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/models"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/repository"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/test"
	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/types"
)

var (
	users = []models.UserModel{
		{
			Username: "TestUser0",
			Email:    "tuo@example.co.uk",
			Password: "",
			FirstName: sql.NullString{
				String: "Test0",
				Valid:  true,
			},
			LastName: sql.NullString{
				String: "User0",
				Valid:  true,
			},
			Role:     types.UserRole,
			Verified: false,
		},
		{
			Username: "TestUser1",
			Email:    "tuo1@example.co.uk",
			Password: "",
			FirstName: sql.NullString{
				String: "Test1",
				Valid:  true,
			},
			LastName: sql.NullString{
				String: "User1",
				Valid:  true,
			},
			Role:     types.UserRole,
			Verified: false,
		},
		{
			Username: "TestUser2",
			Email:    "tuo2@example.co.uk",
			Password: "",
			FirstName: sql.NullString{
				String: "Test2",
				Valid:  true,
			},
			LastName: sql.NullString{
				String: "User2",
				Valid:  true,
			},
			Role:     types.OrganizerRole,
			Verified: false,
		},
	}
)

func TestUserRepository_CreateUser(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	userRepo := repository.NewSQLUserRepository(db)
	t.Run("Create user with empty sql.NullStrings", func(t *testing.T) {
		minimalUser := models.UserModel{
			Role: types.UserRole,
		}

		if err := userRepo.CreateUser(&minimalUser); err != nil {
			t.Errorf("expected no error when creating user with empty sql.NullStrings but got %v", err)
		}

	})
}

func TestUserRepository_KitchenSink(t *testing.T) {
	container, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())

	userRepo := repository.NewSQLUserRepository(db)

	t.Run("Create users", func(t *testing.T) {
		for i := range users {
			if err := userRepo.CreateUser(&users[i]); err != nil {
				t.Errorf("expected no error when creating user but got %v", err)
			}
		}
	})

	t.Run("Update user without loading first", func(t *testing.T) {
		if err := userRepo.UpdateUser(&models.UserModel{Model: models.Model{ID: users[0].ID}, Username: "Updated"}); err == nil {
			t.Error("expected error while attempting to update a user without fully loading the user first")
		}
	})

	t.Run("Update user", func(t *testing.T) {
		users[0].Username = "Updated"
		if err := userRepo.UpdateUser(&users[0]); err != nil {
			t.Errorf("expected no error when updating user but got %v", err)
		}

		loaded, err := userRepo.GetUserByID(users[0].ID)
		if err != nil {
			t.Errorf("expected no error when getting user by id but got %v", err)
		}
		if loaded.Username != "Updated" {
			t.Errorf("expected loaded users username to be updated to 'Updated' but was '%s'", loaded.Username)
		}
	})

	t.Run("Get user by email", func(t *testing.T) {
		loaded, err := userRepo.GetUserByEmail(users[0].Email)
		if err != nil {
			t.Errorf("expected no error when getting user by email but got %v", err)
		}
		if loaded.ID != users[0].ID {
			t.Errorf("expected loaded users id to match '%s' but was '%s'", users[0].ID, loaded.ID)
		}
	})

	// test delete
	t.Run("Delete users", func(t *testing.T) {
		for _, user := range users {
			if err := userRepo.DeleteUser(user.ID); err != nil {
				t.Errorf("expected no error when deleting user but got %v", err)
			}
			dusr, err := userRepo.GetUserByID(user.ID)
			if dusr != nil {
				t.Errorf("expected deleted user to be nil but was %v", dusr)
			}
			if err == nil {
				t.Errorf("expected error when attempting to get user by id after deletion")
			}
		}
	})

	// test insert
}
