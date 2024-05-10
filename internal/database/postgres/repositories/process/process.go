package process

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *postgres.Database) *Repository {
	return &Repository{
		db: database.DB(),
	}
}

func (r *Repository) CreateProcess() (uuid.UUID, error) {
	const op = "database.postgres.CreateProcess"

	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, handleError(op, err)
	}

	_, err = r.db.Exec(`
		INSERT INTO process(id) VALUES ($1);
	`, id)
	if err != nil {
		return uuid.Nil, handleError(op, err)
	}

	return id, nil
}

func (r *Repository) AddOutput(id uuid.UUID, output string, error string) error {
	const op = "database.postgres.AddOutput"

	_, err := r.db.Exec(`
		UPDATE process SET output = output || $1, error = error || $2 WHERE id = $3;
	`, output, error, id)
	if err != nil {
		return handleError(op, err)
	}

	return nil
}

func (r *Repository) FinishProcess(id uuid.UUID, exitCode int) error {
	const op = "database.postgres.FinishProcess"

	var status string
	if exitCode == 0 {
		status = "finished"
	} else {
		status = "error"
	}

	_, err := r.db.Exec(`
		UPDATE process SET status = $1, exit_code = $2 WHERE id = $3;
	`, status, exitCode, id)
	if err != nil {
		return handleError(op, err)
	}

	return nil
}

func (r *Repository) DeleteProcess(id uuid.UUID) error {
	const op = "database.postgres.DeleteProcess"

	err := r.db.QueryRow(`
		DELETE FROM process WHERE id = $1 RETURNING id;
	`, id).Scan(&id)
	if err != nil {
		return handleError(op, err)
	}

	return nil
}

func (r *Repository) Process(id uuid.UUID) (*models.Process, error) {
	const op = "database.postgres.Process"

	var processOutput string
	var processError string
	var processStatus string
	var exitCode int
	err := r.db.QueryRow(`
		SELECT output, error, status, exit_code FROM process WHERE id = $1;
	`, id).Scan(&processOutput, &processError, &processStatus, &exitCode)
	if err != nil {
		return nil, handleError(op, err)
	}

	var status models.ProcessStatus
	switch processStatus {
	case "executing":
		status = models.StatusExecuting
	case "finished":
		status = models.StatusFinished
	case "error":
		status = models.StatusError
	}

	return &models.Process{
		ID:       id,
		Output:   processOutput,
		Error:    processError,
		Status:   status,
		ExitCode: exitCode,
	}, nil
}

func handleError(message string, err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", message, errors.Join(database.ErrProcessNotFound))
		}
		return fmt.Errorf("%s: %w", message, errors.Join(err, database.ErrInternal))
	}

	return nil
}
