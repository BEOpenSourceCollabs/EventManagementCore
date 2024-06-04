package persist

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewDatabase(config DatabaseConfiguration) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode))
	if err == nil {
		err = db.Ping()
	}
	return
}
