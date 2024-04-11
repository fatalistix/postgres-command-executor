package app

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type App struct {
}

func NewApp(log *slog.Logger, port int, timeout, idleTimeout time.Duration) (*App, error) {
	const op = "app.NewApp"

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDbname := os.Getenv("POSTGRES_DBNAME")
	postgresSSLMode := os.Getenv("POSTGRES_SSLMODE")

	var postgresPort uint16
	temp, err := strconv.ParseUint(os.Getenv("POSTGRES_PORT"), 10, 16)
	if err != nil {
		return nil, fmt.Errorf("%s: error parsing port from environment variable: %w", op, err)
	}
	postgresPort = uint16(temp)

	database, err := postgres.NewDatabase(
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDbname,
		postgresSSLMode,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_ = database

	return &App{}, nil
}
