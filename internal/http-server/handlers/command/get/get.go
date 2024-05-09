package get

import (
	"encoding/json"
	"errors"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	ID      int64  `json:"id"`
	Command string `json:"command"`
}

type CommandProvider interface {
	Command(id int64) (models.Command, error)
}

func MakeGetHandlerFunc(log *slog.Logger, provider CommandProvider) http.HandlerFunc {
	const op = "http-server.handlers.command.get.MakeGetHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		pathValueId := r.PathValue("id")
		id, err := strconv.ParseInt(pathValueId, 10, 64)
		if err != nil {
			log.Error("unable to parse id to int64", slog.String("id", pathValueId), slogattr.Err(err))

			http.Error(w, "invalid id", http.StatusBadRequest)

			return
		}

		log.Info("request path value parsed", slog.Int64("id", id))

		command, err := provider.Command(id)
		if err != nil {
			log.Error("error getting command", slog.Int64("id", id), slogattr.Err(err))

			if errors.Is(err, database.ErrCommandNotFound) {
				http.Error(w, "command not found", http.StatusNotFound)
			} else {
				http.Error(w, "error getting command: "+err.Error(), http.StatusBadRequest)
			}

			return
		}

		log.Info("command got", slog.Int64("id", id), slog.String("command", command.Command))

		response := Response{ID: command.ID, Command: command.Command}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err), slog.Any("response", response))

			http.Error(w, "error writing response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
