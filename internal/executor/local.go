package executor

import "os/exec"

type LocalExecutor struct{}

func newLocalExecutor() (*LocalExecutor, error) {
	return &LocalExecutor{}, nil
}

func (local *LocalExecutor) Run(cmd string) (string, error) {
	cmdResult := exec.Command(cmd)
	stdout, err := cmdResult.Output()

	return string(stdout), err
}
