package process

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type GetResponse struct {
	ProcessID uuid.UUID `json:"process_id"`
	Output    string    `json:"output"`
	Error     string    `json:"error"`
	Status    string    `json:"status"`
	ExitCode  int       `json:"exit_code"`
}

type Getter interface {
	Get(id uuid.UUID) (*models.Process, error)
}

func NewGetHandlerFunc(log *slog.Logger, getter Getter) http.HandlerFunc {
	const op = "http-server.handlers.process.get.NewGetHandlerFunc"

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

		process, err := getter.Get(id)
		if err != nil {
			log.Error("error getting process", slogattr.Err(err))

			http.Error(w, "error getting process: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("process got")

		response := GetResponse{
			ProcessID: process.ID,
			Output:    process.Output,
			Error:     process.Error,
			Status:    string(process.Status),
			ExitCode:  process.ExitCode,
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}
	}
}
