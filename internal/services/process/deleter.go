package process

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/domain/wrapper"
	"github.com/fatalistix/postgres-command-executor/internal/lib/syncmap"
	"github.com/google/uuid"
)

type Deleter interface {
	Delete(processID uuid.UUID) error
}

type DeleteService struct {
	sm      *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper]
	deleter Deleter
}

func NewDeleteService(sm *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper], deleter Deleter) *DeleteService {
	return &DeleteService{
		sm: sm,
	}
}

func (d *DeleteService) DeleteProcess(processID uuid.UUID) error {
	const op = "services.process.deleter.DeleteProcess"

	cmd, ok := d.sm.Load(processID)
	if ok {
		killProcess(cmd)
		d.sm.Delete(processID)
	}

	if err := d.deleter.Delete(processID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func killProcess(cmd *wrapper.CmdWrapper) {
	_ = cmd.Cmd.Process.Kill()
	_ = cmd.Stdout.Close()
	_ = cmd.Stderr.Close()
}
