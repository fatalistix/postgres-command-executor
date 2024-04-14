package repositories

import (
	"database/sql"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(database *postgres.Database) (*CommandRepository, error) {
	const op = "database.postgres.NewCommandRepository"

	commandRepository := CommandRepository{db: database.DB()}
	_, err := commandRepository.db.Exec(`
		CREATE TABLE IF NOT EXISTS command (
		    id SERIAL UNIQUE NOT NULL,
			command TEXT UNIQUE NOT NULL,
			PRIMARY KEY(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = commandRepository.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_bash ON command(command);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &commandRepository, nil
}

func (cr *CommandRepository) SaveCommand(command string) (int64, error) {
	const op = "database.postgres.SaveCommand"

	var id int64
	err := cr.db.QueryRow(`
		INSERT INTO command(command) VALUES ($1) RETURNING id;
	`, command).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
