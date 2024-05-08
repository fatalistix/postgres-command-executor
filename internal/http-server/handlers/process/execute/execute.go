package execute

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/lib/http-server/request/header"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type Request struct {
	ID int64 `json:"id"`
}

type Response struct {
	ProcessID uuid.UUID `json:"process_id"`
}

type CommandExecutionStarter interface {
	StartCommandExecution(id int64) (uuid.UUID, error)
}

func MakeExecuteHandlerFunc(log *slog.Logger, executionStarter CommandExecutionStarter) http.HandlerFunc {
	const op = "http-server.handlers.process.NewExecuteHandlerFunc"

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

		log.Info("request body decoded")

		processID, err := executionStarter.StartCommandExecution(request.ID)
		if err != nil {
			log.Error("error starting command execution", slogattr.Err(err))

			http.Error(w, "error starting command execution", http.StatusBadRequest)

			return
		}

		log.Info("command execution started")

		response := Response{ProcessID: processID}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
