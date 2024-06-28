package test

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabaseConfiguration struct {
	RootRelativePath string
}

func NewTestDatabaseWithContainer(config TestDatabaseConfiguration) (*postgres.PostgresContainer, *sql.DB, error) {
	postgresContainer, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("postgres:latest"),
		testcontainers.WithEnv(map[string]string{
			"DATABASE_USER":     "postgres",
			"DATABASE_NAME":     "event-mgmt-db",
			"DATABASE_PORT":     "5432",
			"DATABASE_SSL_MODE": "disable",
		}),
		postgres.WithInitScripts(filepath.Join(config.RootRelativePath, "init.sql"), filepath.Join(config.RootRelativePath, "migrations/000001_initial.up.sql"), filepath.Join(config.RootRelativePath, "migrations/000002_add_google_id_avatar_to_users.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
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
