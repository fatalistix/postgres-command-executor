package process

import (
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

type CommandGetter interface {
	Get(id int64) (models.Command, error)
}

type Provider interface {
	Create() (uuid.UUID, error)
	AddOutput(processID uuid.UUID, output string, error string) error
	Finish(processID uuid.UUID, exitCode int) error
	Delete(processID uuid.UUID) error
}

type ExecuteService struct {
	commandGetter CommandGetter
	provider      Provider
	sm            *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper]
}

func NewProcessExecutor(getter CommandGetter, provider Provider, sm *syncmap.SyncMap[uuid.UUID, *wrapper.CmdWrapper]) *ExecuteService {
	return &ExecuteService{
		commandGetter: getter,
		provider:      provider,
		sm:            sm,
	}
}

type readResult struct {
	Result string
	Err    error
}

func (e *ProcessExecutor) StartCommandExecution(commandID int64) (uuid.UUID, error) {
	const op = "services.process.executor.ExecuteCommand"

	command, err := e.commandGetter.GetCommand(commandID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	processID, err := e.processProvider.CreateProcess()
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

	go e.listenStdoutAndStderr(processID, stdoutCh, stderrCh)

	e.sm.Store(processID, &wrapper.CmdWrapper{
		Cmd:    cmd,
		Stdout: stdout,
		Stderr: stderr,
	})

	return processID, nil
}

func (e *ProcessExecutor) listenStdoutAndStderr(processID uuid.UUID, stdoutCh chan readResult, stderrCh chan readResult) {
	var stdoutErr error
	var stderrErr error
	for {
		select {
		case result := <-stdoutCh:
			_ = e.processProvider.AddOutput(processID, result.Result, "")
			if result.Err != nil {
				_ = e.processProvider.FinishProcess(processID, 1)
			}
		case result := <-stderrCh:
			_ = e.processProvider.AddOutput(processID, "", result.Result)
			if result.Err != nil {
				_ = e.processProvider.FinishProcess(processID, 1)
			}
		}
		if stdoutErr != nil && stderrErr != nil {
			if stdoutErr == io.EOF && stderrErr == io.EOF {
				_ = e.processProvider.FinishProcess(processID, 0)
				return
			}
		}
	}
}

func readInLoop(reader io.Reader, resultCh chan<- readResult) {
	buffer := make([]byte, bufferSize)

	for {
		n, err := reader.Read(buffer)
		resultCh <- readResult{Result: string(buffer[:n]), Err: err}
		if err != nil {
			return
		}
	}
}
