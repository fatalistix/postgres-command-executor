package save

import (
	"encoding/json"
	"errors"
	"github.com/fatalistix/postgres-command-executor/internal/database"
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

func MakeSaveHandlerFunc(log *slog.Logger, saver CommandSaver) http.HandlerFunc {
	const op = "http-server.handlers.command.save.MakeSaveHandlerFunc"

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

		id, err := saver.SaveCommand(request.Command)
		if err != nil {
			log.Error("error saving command", slog.Any("request", request), slogattr.Err(err))

			if errors.Is(err, database.ErrCommandExists) {
				http.Error(w, "command already exists", http.StatusConflict)
			} else {
				http.Error(w, "error saving command", http.StatusInternalServerError)
			}

			return
		}

		log.Info("new command saved", slog.Int64("id", id), slog.String("command", request.Command))

		response := Response{ID: id}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err), slog.Any("response", response))

			http.Error(w, "error writing response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
