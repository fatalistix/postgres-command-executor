package command

import (
	"database/sql"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *postgres.Database) (*Repository, error) {
	const op = "database.postgres.NewCommandRepository"

	commandRepository := Repository{db: database.DB()}
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

func (cr *Repository) SaveCommand(command string) (int64, error) {
	const op = "database.postgres.repositories.SaveCommand"

	var id int64
	err := cr.db.QueryRow(`
		INSERT INTO command(command) VALUES ($1) RETURNING id;
	`, command).Scan(&id)
	if err != nil {
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

func (cr *Repository) GetCommands() ([]models.Command, error) {
	const op = "database.postgres.repositories.GetCommands"

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

func (cr *Repository) GetCommand(id int64) (models.Command, error) {
	const op = "database.postgres.repositories.GetCommand"

	var command models.Command

	err := cr.db.QueryRow(`
		SELECT id, command FROM command WHERE id = $1
	`, id).Scan(&command.ID, &command.Command)
	if err != nil {
		return models.Command{}, fmt.Errorf("%s: %w", op, err)
	}

	return command, nil
}
