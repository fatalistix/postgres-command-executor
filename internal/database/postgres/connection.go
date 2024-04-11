package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(host string, port uint16, user string, password string, dbname string, sslMode string) (*Database, error) {
	const op = "database.postgres.NewDatabase"

	psqlCredentials := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode,
	)

	db, err := sql.Open("postgres", psqlCredentials)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// verify that data source name is valid (according to godoc of `sql.Open`)
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Database{db: db}, nil
}
