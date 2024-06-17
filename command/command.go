package command

import (
	"os"
	"os/exec"
)

type Command struct {
	cmd *exec.Cmd
}

func New(name string, arg ...string) *Command {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &Command{
		cmd: cmd,
	}
}

func (c *Command) Run() error {
	return c.cmd.Run()
}
