package executor

import (
	"log"
	"os/exec"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/common"
)

type LocalExecutor struct{}

func newLocalExecutor() (*LocalExecutor, error) {
	return &LocalExecutor{}, nil
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
			log.Panic(err)
		}

		log.Print(str)
	}
}
