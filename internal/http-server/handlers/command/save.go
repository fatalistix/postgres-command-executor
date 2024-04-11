package command

import "log/slog"

type Request struct {
	Bash string `json:"bash"`
}

type CommandSaver interface {
	SaveCommand
}

func NewSaveHandlerFunc(log *slog.Logger)
