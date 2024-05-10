package models

import "github.com/google/uuid"

type ProcessStatus string

const (
	StatusExecuting ProcessStatus = "executing"
	StatusFinished  ProcessStatus = "finished"
	StatusError     ProcessStatus = "error"
)

type Process struct {
	ID       uuid.UUID
	Output   string
	Error    string
	Status   ProcessStatus
	ExitCode int
}
