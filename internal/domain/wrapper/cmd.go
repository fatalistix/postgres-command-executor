package wrapper

import (
	"io"
	"os/exec"
)

type CmdWrapper struct {
	Cmd    *exec.Cmd
	Stdout io.Closer
	Stderr io.Closer
}
