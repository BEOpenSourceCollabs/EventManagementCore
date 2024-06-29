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
		for _, user := range users {
			if err := userRepo.CreateUser(&user); err != nil {
				t.Errorf("expected no error when creating user but got %v", err)
			}
		}
	})

	// TODO:
	// test update
	// test delete
	// test insert
}
