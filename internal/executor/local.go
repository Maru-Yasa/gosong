package executor

import (
	"os"
	"os/exec"
	"strings"
)

type LocalExecutor struct {
	name string
}

func (local *LocalExecutor) GetName() string {
	return local.name
}

func newLocalExecutor(name string) (*LocalExecutor, error) {
	return &LocalExecutor{
		name: name,
	}, nil
}

func (local *LocalExecutor) RunRaw(cmd string) (string, error) {
	return "local executor does not support RunRaw yet", nil
}

func (local *LocalExecutor) Run(cmd string, cwd string) error {
	parts := strings.Fields(cmd)
	program := parts[0]
	args := parts[1:]
	cmdResult := exec.Command(program, args...)
	if cwd != "" {
		cmdResult.Dir = cwd
	}

	cmdResult.Stdout = os.Stdout
	cmdResult.Stderr = os.Stderr

	err := cmdResult.Run()
	return err
}
