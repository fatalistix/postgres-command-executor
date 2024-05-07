package command

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
	"strconv"
)

type GetResponse struct {
	ID      int64  `json:"id"`
	Command string `json:"command"`
}

type Getter interface {
	GetCommand(id int64) (models.Command, error)
}

func NewGetHandlerFunc(log *slog.Logger, getter Getter) http.HandlerFunc {
	const op = "http-server.handlers.command.NewGetHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		pathValueId := r.PathValue("id")
		id, err := strconv.ParseInt(pathValueId, 10, 64)
		if err != nil {
			log.Error("unable to parse id '"+pathValueId+"' to int64", slogattr.Err(err))

			http.Error(w, "unable to parse id to int64", http.StatusBadRequest)

			return
		}

		log.Info("request path value parsed", slog.Any("id", id))

		command, err := getter.Get(id)
		if err != nil {
			log.Error("error getting command", slogattr.Err(err))

			http.Error(w, "error getting command: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("command got")

		response := GetResponse{ID: command.ID, Command: command.Command}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
