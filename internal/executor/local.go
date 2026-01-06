package executor

import (
	"os"
	"os/exec"
)

type LocalExecutor struct {
	name  string
	shell string
}

func (local *LocalExecutor) GetName() string {
	return local.name
}

func newLocalExecutor(name string) (*LocalExecutor, error) {
	return &LocalExecutor{
		name:  name,
		shell: "/bin/sh", // default shell
	}, nil
}

func (local *LocalExecutor) RunRaw(cmd string) (string, error) {
	cmdResult := exec.Command(local.shell, "-c", cmd)

	output, err := cmdResult.CombinedOutput()
	return string(output), err
}

func (local *LocalExecutor) Run(cmd string, cwd string) error {
	cmdResult := exec.Command(local.shell, "-c", cmd)
	if cwd != "" {
		cmdResult.Dir = cwd
	}

	cmdResult.Stdout = os.Stdout
	cmdResult.Stderr = os.Stderr

	err := cmdResult.Run()
	return err
}
