package execute

import (
	"errors"
	"fmt"
	"github.com/fatalistix/postgres-command-executor/internal/domain/models"
	"github.com/fatalistix/postgres-command-executor/internal/domain/wrapper"
	"github.com/fatalistix/postgres-command-executor/internal/lib/syncmap"
	"github.com/google/uuid"
	"io"
	"os/exec"
)

const (
	bufferSize = 1024
)

type CommandProvider interface {
	Command(id int64) (models.Command, error)
}

type ProcessProvider interface {
	CreateProcess() (uuid.UUID, error)
	AddOutput(processID uuid.UUID, output string, error string) error
	FinishProcess(processID uuid.UUID, exitCode int) error
	DeleteProcess(processID uuid.UUID) error
}

type Service struct {
	commandProvider CommandProvider
	processProvider ProcessProvider
	sm              *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper]
}

func NewService(
	commandProvider CommandProvider,
	processProvider ProcessProvider,
	sm *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper],
) *Service {
	return &Service{
		commandProvider: commandProvider,
		processProvider: processProvider,
		sm:              sm,
	}
}

type readResult struct {
	Result string
	Err    error
}

func (s *Service) StartCommandExecution(commandID int64) (uuid.UUID, error) {
	const op = "services.process.executor.ExecuteCommand"

	command, err := s.commandProvider.Command(commandID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	processID, err := s.processProvider.CreateProcess()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	cmd := exec.Command("/bin/bash", "-c", command.Command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = cmd.Start()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stdoutCh := make(chan readResult)
	go readInLoop(stdout, stdoutCh)

	stderrCh := make(chan readResult)
	go readInLoop(stderr, stderrCh)

	go s.listenStdoutAndStderr(processID, stdoutCh, stderrCh)

	s.sm.Store(processID, &wrapper.CmdWrapper{
		Cmd:    cmd,
		Stdout: stdout,
		Stderr: stderr,
	})

	return processID, nil
}

func (s *Service) listenStdoutAndStderr(processID uuid.UUID, stdoutCh chan readResult, stderrCh chan readResult) {
	var stdoutErr error
	var stderrErr error
	for {
		select {
		case result := <-stdoutCh:
			_ = s.processProvider.AddOutput(processID, result.Result, "")
			if result.Err != nil {
				stdoutErr = result.Err
			}
		case result := <-stderrCh:
			_ = s.processProvider.AddOutput(processID, "", result.Result)
			if result.Err != nil {
				stderrErr = result.Err
			}
		}
		if stdoutErr != nil && stderrErr != nil {
			if errors.Is(stdoutErr, io.EOF) && errors.Is(stderrErr, io.EOF) {
				_ = s.processProvider.FinishProcess(processID, 0)
			} else {
				_ = s.processProvider.FinishProcess(processID, 1)
			}
			return
		}
	}
}

func readInLoop(reader io.Reader, resultCh chan<- readResult) {
	const op = "services.process.executor.readInLoop"

	buffer := make([]byte, bufferSize)

	for {
		n, err := reader.Read(buffer)
		resultCh <- readResult{
			Result: string(buffer[:n]),
			Err:    fmt.Errorf("%s: %w", op, err),
		}
		if err != nil {
			return
		}
	}
}
