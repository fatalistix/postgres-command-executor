package delete

import (
	"errors"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
	"strconv"
)

type CommandDeleter interface {
	DeleteCommand(id int64) error
}

func MakeDeleteHandlerFunc(log *slog.Logger, deleter CommandDeleter) http.HandlerFunc {
	const op = "http-server.handlers.command.delete.MakeDeleteHandlerFunc"

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

		err = deleter.DeleteCommand(id)
		if err != nil {
			log.Error("error deleting command", slog.Int64("id", id), slogattr.Err(err))

			if errors.Is(err, database.ErrCommandNotFound) {
				http.Error(w, "command not found", http.StatusNotFound)
			} else {
				http.Error(w, "error deleting command", http.StatusInternalServerError)
			}

			return
		}

		log.Info("command deleted", slog.Int64("id", id))

		w.WriteHeader(http.StatusOK)

		log.Info("response is written")
	}
}
