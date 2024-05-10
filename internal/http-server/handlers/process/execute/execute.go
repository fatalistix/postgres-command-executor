package execute

import (
	"encoding/json"
	"errors"
	"github.com/fatalistix/postgres-command-executor/internal/database"
	"github.com/fatalistix/postgres-command-executor/internal/lib/http-server/request/header"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Request struct {
	CommandID int64 `json:"command_id"`
}

type Response struct {
	ProcessID uuid.UUID `json:"process_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=CommandExecutionStarter
type CommandExecutionStarter interface {
	StartCommandExecution(id int64) (uuid.UUID, error)
}

func MakeExecuteHandlerFunc(log *slog.Logger, executionStarter CommandExecutionStarter) http.HandlerFunc {
	const op = "http-server.handlers.process.MakeExecuteHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if !header.HasApplicationJson(r) {
			log.Error("no 'application/json' header found")

			http.Error(w, "no 'application/json' header found", http.StatusUnsupportedMediaType)

			return
		}

		var request Request

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			log.Error("unable to decode request's body", slogattr.Err(err))

			http.Error(w, "unable to decode request's body", http.StatusBadRequest)

			return
		}

		log.Info("request body decoded", slog.Any("request", request))

		processID, err := executionStarter.StartCommandExecution(request.CommandID)
		if err != nil {
			log.Error("error starting command execution", slogattr.Err(err), slog.Any("request", request))

			if errors.Is(err, database.ErrCommandNotFound) {
				http.Error(w, "command not found", http.StatusNotFound)
			} else {
				http.Error(w, "error starting command execution", http.StatusInternalServerError)
			}

			return
		}

		log.Info("command execution started", slog.Any("process_id", processID))

		response := Response{ProcessID: processID}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err), slog.Any("response", response))

			http.Error(w, "error writing response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
