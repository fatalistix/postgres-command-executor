package delete

import (
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

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
			log.Error("unable to parse id '"+pathValueId+"' to UUID", slogattr.Err(err))

			http.Error(w, "unable to parse id to UUID", http.StatusBadRequest)

			return
		}

		log.Info("process id parsed", slog.Any("id", id))

		err = deleter.DeleteProcess(id)
		if err != nil {
			log.Error("error deleting process", slogattr.Err(err))

			http.Error(w, "error deleting process: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("process deleted")

		w.WriteHeader(http.StatusOK)

		log.Info("response written")
	}
}
