package get

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Response struct {
	ID       uuid.UUID `json:"id"`
	Output   string    `json:"output"`
	Error    string    `json:"error"`
	Status   string    `json:"status"`
	ExitCode int       `json:"exit_code"`
}

type ProcessProvider interface {
	Process(id uuid.UUID) (*models.Process, error)
}

func MakeGetHandlerFunc(log *slog.Logger, provider ProcessProvider) http.HandlerFunc {
	const op = "http-server.handlers.process.get.MakeGetHandlerFunc"

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

		process, err := provider.Process(id)
		if err != nil {
			log.Error("error getting process", slogattr.Err(err))

			http.Error(w, "error getting process: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("process got")

		response := Response{
			ID:       process.ID,
			Output:   process.Output,
			Error:    process.Error,
			Status:   string(process.Status),
			ExitCode: process.ExitCode,
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}
	}
}
