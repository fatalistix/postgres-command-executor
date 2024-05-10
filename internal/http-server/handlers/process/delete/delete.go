package delete

import (
	"errors"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=ProcessDeleter
type ProcessDeleter interface {
	DeleteProcess(id uuid.UUID) error
}

func MakeDeleteHandlerFunc(log *slog.Logger, deleter ProcessDeleter) http.HandlerFunc {
	const op = "http-server.handlers.process.delete.MakeDeleteHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		pathValueId := r.PathValue("id")
		id, err := uuid.Parse(pathValueId)
		if err != nil {
			log.Error("unable to parse id to UUID", slogattr.Err(err), slog.String("id", pathValueId))

			http.Error(w, "invalid id", http.StatusBadRequest)

			return
		}

		log.Info("process id parsed", slog.Any("id", id))

		err = deleter.DeleteProcess(id)
		if err != nil {
			log.Error("error deleting process", slogattr.Err(err), slog.Any("id", id))

			if errors.Is(err, database.ErrProcessNotFound) {
				http.Error(w, "process not found", http.StatusNotFound)
			} else {
				http.Error(w, "error deleting process", http.StatusInternalServerError)
			}

			return
		}

		log.Info("process deleted")

		w.WriteHeader(http.StatusOK)

		log.Info("response written")
	}
}
