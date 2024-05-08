package process

import (
	"database/sql"
	"fmt"
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
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Exec(`
		INSERT INTO process(id) VALUES ($1);
	`, id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) AddOutput(id uuid.UUID, output string, error string) error {
	const op = "database.postgres.AddOutput"

	_, err := r.db.Exec(`
		UPDATE process SET output = output || $1, error = error || $2 WHERE id = $3;
	`, output, error, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) FinishProcess(id uuid.UUID, exitCode int) error {
	const op = "database.postgres.FinishProcess"

	var status string
	if exitCode != 0 {
		status = "error"
	} else {
		status = "finished"
	}

	_, err := r.db.Exec(`
		UPDATE process SET status = $1, exit_code = $2 WHERE id = $3;
	`, status, exitCode, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) DeleteProcess(id uuid.UUID) error {
	const op = "database.postgres.DeleteProcess"

	_, err := r.db.Exec(`
		DELETE FROM process WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var status models.ProcessStatus
	switch processStatus {
	case "executing":
		status = models.Executing
	case "finished":
		status = models.Finished
	case "error":
		status = models.Error
	}

	return &models.Process{
		ID:       id,
		Output:   processOutput,
		Error:    processError,
		Status:   status,
		ExitCode: exitCode,
	}, nil
}
