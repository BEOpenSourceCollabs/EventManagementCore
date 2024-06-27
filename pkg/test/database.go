package test

import (
	"context"
	"database/sql"
	"time"

	"github.com/BEOpenSourceCollabs/EventManagementCore/pkg/persist"
	_ "github.com/mattn/go-sqlite3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewTestDatabaseWithContainer(config persist.DatabaseConfiguration) (*postgres.PostgresContainer, *sql.DB, error) {
	postgresContainer, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:latest"),
		// postgres.WithInitScripts("../../init.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	// Run any migrations on the database
	_, _, err = postgresContainer.Exec(context.Background(), []string{"psql", "-U", config.User, "-d", config.Database, "-f", "../../init.up.sql"})
	if err != nil {
		return nil, nil, err
	}
	_, _, err = postgresContainer.Exec(context.Background(), []string{"psql", "-U", config.User, "-d", config.Database, "-f", "../../migrations/000001_initial.up.sql"})
	if err != nil {
		return nil, nil, err
	}
	_, _, err = postgresContainer.Exec(context.Background(), []string{"psql", "-U", config.User, "-d", config.Database, "-f", "../../migrations/000002_add_google_id_avatar_to_users.up.sql"})
	if err != nil {
		return nil, nil, err
	}

	conn, _ := postgresContainer.ConnectionString(context.Background(), "sslmode=disable")
	db, err := sql.Open("postgres", conn)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, nil, err
	}

	return postgresContainer, db, nil
}
