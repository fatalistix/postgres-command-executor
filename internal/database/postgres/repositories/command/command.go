package command

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *postgres.Database) *Repository {
	return &Repository{
		db: database.DB(),
	}
}

func (cr *Repository) SaveCommand(command string) (int64, error) {
	const op = "database.postgres.repositories.SaveCommand"

	var id int64
	err := cr.db.QueryRow(`
		INSERT INTO command(command) VALUES ($1) RETURNING id;
	`, command).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return 0, fmt.Errorf("%s: %w", op, database.ErrCommandExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (cr *Repository) DeleteCommand(id int64) error {
	const op = "database.postgres.repositories.DeleteCommand"

	_, err := cr.db.Exec(`
		DELETE FROM command WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (cr *Repository) Commands() ([]models.Command, error) {
	const op = "database.postgres.repositories.Commands"

	rows, err := cr.db.Query(`
		SELECT id, command FROM command;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = rows.Close()
	}()

	commands := make([]models.Command, 0)

	for rows.Next() {
		var command models.Command
		if err := rows.Scan(&command.ID, &command.Command); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		commands = append(commands, command)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return commands, nil
}

func (cr *Repository) Command(id int64) (models.Command, error) {
	const op = "database.postgres.repositories.Command"

	var command models.Command

	err := cr.db.QueryRow(`
		SELECT id, command FROM command WHERE id = $1
	`, id).Scan(&command.ID, &command.Command)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Command{}, fmt.Errorf("%s: %w", op, database.ErrCommandNotFound)
		}
		return models.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	return command, nil
}
