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
		_ = s.processProvider.DeleteProcess(processID)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		_ = s.processProvider.DeleteProcess(processID)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = cmd.Start()
	if err != nil {
		_ = stderr.Close()
		_ = stdout.Close()
		_ = s.processProvider.DeleteProcess(processID)
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	stdoutCh := make(chan readResult)
	go readInLoop(stdout, stdoutCh)

	stderrCh := make(chan readResult)
	go readInLoop(stderr, stderrCh)

	go s.listenStdoutAndStderr(processID, cmd, stdoutCh, stderrCh)

	s.sm.Store(processID, &wrapper.CmdWrapper{
		Cmd:    cmd,
		Stdout: stdout,
		Stderr: stderr,
	})

	return processID, nil
}

func (s *Service) listenStdoutAndStderr(processID uuid.UUID, cmd *exec.Cmd, stdoutCh chan readResult, stderrCh chan readResult) {
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
			err := cmd.Wait()
			if err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) {
					_ = s.processProvider.FinishProcess(processID, exitError.ExitCode())
				} else {
					_ = s.processProvider.FinishProcess(processID, 1)
				}
			} else {
				_ = s.processProvider.FinishProcess(processID, 0)
			}
			s.sm.Delete(processID)
			return
		}
	}
}

func readInLoop(readCloser io.ReadCloser, resultCh chan<- readResult) {
	buffer := make([]byte, bufferSize)

	for {
		n, err := readCloser.Read(buffer)
		resultCh <- readResult{
			Result: string(buffer[:n]),
			Err:    err,
		}
		if err != nil {
			_ = readCloser.Close()
			return
		}
	}
}
