package save

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/lib/http-server/request/header"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
)

type Request struct {
	Command string `json:"command"`
}

type Response struct {
	ID int64 `json:"id"`
}

type CommandSaver interface {
	SaveCommand(command string) (int64, error)
}

func NewSaveHandlerFunc(log *slog.Logger, saver CommandSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.command.save.NewSaveHandlerFunc"

		log = log.With(
			slog.String("op", op),
		)

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

		id, err := saver.SaveCommand(request.Command)
		if err != nil {
			log.Error("error saving command", slogattr.Err(err))

			http.Error(w, "error saving command: "+err.Error(), http.StatusBadRequest)

			return
		}

		response := Response{ID: id}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}
	}
}
