package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

type CommandRepository struct {
	db *sql.DB
}

func NewCommandRepository(database *Database) (*CommandRepository, error) {
	const op = "database.postgres.NewCommandRepository"

	commandRepository := CommandRepository{db: database.db}
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

func (cr *CommandRepository) SaveCommand(command string) error {
	const op = "database.postgres.SaveCommand"

	stmt, err := cr.db.Prepare(`
		INSERT INTO command(command) values(?);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(stmt, command)
	if err != nil {
		var pqErr *pq.Error
		errors.As(err, &pqErr)
		fmt.Println(pqErr.Code.Name())
	}

	return nil
}
