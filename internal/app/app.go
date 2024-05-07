package app

import (
	"context"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/config"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories"
	commanddelete "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/delete"
	commandget "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/get"
	commandlist "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/list"
	commandsave "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save"
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

	commandRepository, err := repositories.NewCommandRepository(database)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /commands", commandsave.NewSaveHandlerFunc(log, commandRepository))
	mux.HandleFunc("DELETE /command/{id}", commanddelete.NewDeleteHandlerFunc(log, commandRepository))
	mux.HandleFunc("GET /commands", commandlist.NewListHandlerFunc(log, commandRepository))
	mux.HandleFunc("GET /command/{id}", commandget.NewGetHandlerFunc(log, commandRepository))

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