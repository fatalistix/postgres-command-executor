package delete

import (
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/domain/wrapper"
	"github.com/fatalistix/postgres-command-executor/internal/lib/syncmap"
	"github.com/google/uuid"
)

type ProcessDeleter interface {
	DeleteProcess(processID uuid.UUID) error
}

type Service struct {
	sm      *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper]
	deleter ProcessDeleter
}

func NewService(sm *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper], deleter ProcessDeleter) *Service {
	return &Service{
		sm:      sm,
		deleter: deleter,
	}
}

func (s *Service) DeleteProcess(processID uuid.UUID) error {
	const op = "services.process.delete.DeleteProcess"

	cmd, ok := s.sm.Load(processID)
	if ok {
		killProcess(cmd)
		s.sm.Delete(processID)
	}

	if err := s.deleter.DeleteProcess(processID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func killProcess(cmd *wrapper.CmdWrapper) {
	_ = cmd.Cmd.Process.Kill()
	_ = cmd.Stdout.Close()
	_ = cmd.Stderr.Close()
}
