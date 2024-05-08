package command

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
)

type ListResponse struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	ID      int64  `json:"id"`
	Command string `json:"command"`
}

type ListGetter interface {
	GetList() ([]models.Command, error)
}

func NewListHandlerFunc(log *slog.Logger, listGetter ListGetter) http.HandlerFunc {
	const op = "http-server.handlers.command.NewListHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request received")

		commands, err := listGetter.GetList()
		if err != nil {
			log.Error("error getting commands", slogattr.Err(err))

			http.Error(w, "error getting commands: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("commands got")

		response := ListResponse{Commands: make([]Command, 0, len(commands))}
		for _, c := range commands {
			response.Commands = append(response.Commands, Command{c.ID, c.Command})
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err))

			http.Error(w, "error encoding response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}