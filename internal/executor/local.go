package executor

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Maru-Yasa/gosong/pkg/logger"

	"github.com/Maru-Yasa/gosong/internal/common"
)

type LocalExecutor struct {
	name string
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
	stdout, err := cmdResult.Output()

	return string(stdout), err
}

func (local *LocalExecutor) RunTask(task *common.Task) {
	for _, step := range task.Steps {
		str, err := local.Run(step.Command)

		if err != nil {
			logger.Error(fmt.Sprint("Local command failed: ", err), local.name)
			panic(err)
		}

		logger.Info(fmt.Sprint("Local command output: ", strings.TrimSpace(str)), local.name)
	}
}
