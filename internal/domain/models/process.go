package models

import "github.com/google/uuid"

type ProcessStatus string

const (
	Executing ProcessStatus = "executing"
	Finished  ProcessStatus = "finished"
	Error     ProcessStatus = "error"
)

type Process struct {
	ID       uuid.UUID
	Output   string
	Error    string
	Status   ProcessStatus
	ExitCode int
}
