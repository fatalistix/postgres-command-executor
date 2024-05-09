package list

import (
	"encoding/json"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
)

type Response struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	ID      int64  `json:"id"`
	Command string `json:"command"`
}

type CommandProvider interface {
	Commands() ([]models.Command, error)
}

func MakeListHandlerFunc(log *slog.Logger, provider CommandProvider) http.HandlerFunc {
	const op = "http-server.handlers.command.list.MakeListHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request received")

		commands, err := provider.Commands()
		if err != nil {
			log.Error("error getting commands", slogattr.Err(err))

			http.Error(w, "error getting commands", http.StatusInternalServerError)

			return
		}

		log.Info("commands got")

		response := Response{Commands: make([]Command, 0, len(commands))}
		for _, c := range commands {
			response.Commands = append(response.Commands, Command{c.ID, c.Command})
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			log.Error("error encoding response", slogattr.Err(err), slog.Any("response", response))

			http.Error(w, "error writing response", http.StatusInternalServerError)

			return
		}

		log.Info("response is written")
	}
}
