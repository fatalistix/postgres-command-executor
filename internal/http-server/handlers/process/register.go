package process

import (
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories/command"
	"github.com/fatalistix/postgres-command-executor/internal/database/postgres/repositories/process"
	"github.com/fatalistix/postgres-command-executor/internal/domain/wrapper"
	deletehandler "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/delete"
	executehandler "github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/execute"
	"github.com/fatalistix/postgres-command-executor/internal/http-server/handlers/process/get"
	"github.com/fatalistix/postgres-command-executor/internal/lib/syncmap"
	deleteservice "github.com/fatalistix/postgres-command-executor/internal/services/process/delete"
	executeservice "github.com/fatalistix/postgres-command-executor/internal/services/process/execute"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

func RegisterHandlers(
	router *http.ServeMux,
	log *slog.Logger,
	processRepository *process.Repository,
	commandRepository *command.Repository,
) {
	sm := syncmap.NewSyncMap[uuid.UUID, *wrapper.CmdWrapper]()
	executeService := executeservice.NewService(commandRepository, processRepository, sm)
	deleteService := deleteservice.NewService(processRepository, sm)

	router.Handle("DELETE /process/{id}", deletehandler.MakeDeleteHandlerFunc(log, deleteService))
	router.Handle("POST /processes", executehandler.MakeExecuteHandlerFunc(log, executeService))
	router.Handle("GET /process/{id}", get.MakeGetHandlerFunc(log, processRepository))
}
