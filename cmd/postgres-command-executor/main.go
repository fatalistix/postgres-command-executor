package main

import (
	"context"
	"github.com/fatalistix/postgres-command-executor/internal/app"
	"github.com/fatalistix/postgres-command-executor/internal/config"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := setupLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error("failed to load environment variables", slogattr.Err(err))
		os.Exit(1)
	}

	pathToConfig := os.Getenv("CONFIG_PATH")

	cfg := config.MustLoadConfig(pathToConfig)

	application, err := app.NewApp(log, cfg)
	if err != nil {
		log.Error("failed to init application", slogattr.Err(err))
		os.Exit(1)
	}

	log.Info("starting application", slog.String("address", cfg.HTTPServer.Address))

	go func() {
		if err := application.Run(); err != nil {
			log.Error("server stopped with error", slogattr.Err(err))
		}
	}()

	log.Info("application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := application.Stop(ctx); err != nil {
		log.Error("failed to stop application", slogattr.Err(err))

		return
	}

	log.Info("application stopped")
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
}
