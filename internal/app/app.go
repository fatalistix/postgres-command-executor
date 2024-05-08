package app

import (
	"context"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/config"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories/command"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories/process"
	commandhandlers "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command"
	processhandlers "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process"
	"log/slog"
	"net/http"
)

type App struct {
	srv      *http.Server
	database *postgres.Database
}

func NewApp(log *slog.Logger, cfg config.Config) (*App, error) {
	const op = "app.NewApp"

	database, err := postgres.NewDatabase(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	commandRepository := command.NewRepository(database)
	processRepository := process.NewRepository(database)

	mux := http.NewServeMux()

	commandhandlers.RegisterHandlers(mux, log, commandRepository)
	processhandlers.RegisterHandlers(mux, log, processRepository, commandRepository)

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{srv: &srv, database: database}, nil
}

func (a *App) Run() error {
	const op = "app.Run"

	if err := a.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.Stop"

	defer func() {
		_ = a.database.DB().Close()
	}()

	if err := a.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
