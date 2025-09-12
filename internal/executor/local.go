package executor

import (
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

func (local *LocalExecutor) Run(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	program := parts[0]
	args := parts[1:]
	cmdResult := exec.Command(program, args...)
	output, err := cmdResult.CombinedOutput()
	return string(output), err
}
