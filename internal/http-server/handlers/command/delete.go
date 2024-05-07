package command

import (
	slogattr "github.com/fatalistix/postgres-command-executor/internal/lib/log/slog/attr"
	"log/slog"
	"net/http"
	"strconv"
)

type Deleter interface {
	Delete(id int64) error
}

func NewDeleteHandlerFunc(log *slog.Logger, deleter Deleter) http.HandlerFunc {
	const op = "http-server.handlers.command.NewDeleteHandlerFunc"

	log = log.With(
		slog.String("op", op),
	)

	return func(w http.ResponseWriter, r *http.Request) {
		pathValueId := r.PathValue("id")
		id, err := strconv.ParseInt(pathValueId, 10, 64)
		if err != nil {
			log.Error("unable to parse id '"+pathValueId+"' to int64", slogattr.Err(err))

			log.Error("unable to parse id to int64", slogattr.Err(err))

			http.Error(w, "unable to parse id to int64", http.StatusBadRequest)

			return
		}

		log.Info("request path value parsed")

		err = deleter.Delete(id)
		if err != nil {
			log.Error("error deleting command", slogattr.Err(err))

			http.Error(w, "error saving command: "+err.Error(), http.StatusBadRequest)

			return
		}

		log.Info("command deleted")

		w.WriteHeader(http.StatusOK)

		log.Info("response written")
	}
}
