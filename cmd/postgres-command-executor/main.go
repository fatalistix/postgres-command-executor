package main

import (
	"github.com/fatalistix/postgres-command-executor/internal/app"
	"github.com/fatalistix/postgres-command-executor/internal/config"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

func main() {
	log := setupLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error("failed to load environment variables", slogattr.Err(err))
		os.Exit(1)
	}

	cfg := config.MustLoadConfig()

	application, err := app.NewApp(log, cfg.HTTPServer.Port, cfg.HTTPServer.Timeout, cfg.HTTPServer.IdleTimeout)
	if err != nil {
		log.Error("failed to init application", slogattr.Err(err))
		os.Exit(1)
	}

	_ = application
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
}
