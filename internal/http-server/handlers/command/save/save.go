package command

import (
	"log/slog"
	"net/http"
)

type Request struct {
	Bash string `json:"bash"`
}

type Saver interface {
	SaveCommand()
}

func NewSaveHandlerFunc(log *slog.Logger, saver Saver) *http.HandlerFunc {
	return func() {

	}
}
