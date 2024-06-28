package repository_test

import (
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

func TestUserRepository_KitchenSink(t *testing.T) {
	_, db, err := test.NewTestDatabaseWithContainer(test.TestDatabaseConfiguration{
		RootRelativePath: "../../",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Clean up the container
	// defer func() {
	// 	if err := container.Terminate(context.Background()); err != nil {
	// 		log.Fatalf("failed to terminate container: %s", err)
	// 	}
	// }()
	userRepo := repository.NewSQLUserRepository(db)

	t.Run("Create users", func(t *testing.T) {
		for _, user := range users {
			if err := userRepo.CreateUser(&user); err != nil {
				t.Errorf("expected no error when creating user but got %v", err)
			}
		}
	})
}
