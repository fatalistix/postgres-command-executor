package command

import (
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories/command"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/delete"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/get"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/list"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/command/save"
	"log/slog"
	"net/http"
)

func RegisterHandlers(router *http.ServeMux, log *slog.Logger, repository *command.Repository) {
	router.Handle("DELETE /command/{id}", delete.MakeDeleteHandlerFunc(log, repository))
	router.Handle("GET /command/{id}", get.MakeGetHandlerFunc(log, repository))
	router.Handle("GET /commands", list.MakeListHandlerFunc(log, repository))
	router.Handle("POST /commands", save.MakeSaveHandlerFunc(log, repository))
}
