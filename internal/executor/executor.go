package executor

import (
	"fmt"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

type Executor interface {
	Run(cmd string) (string, error)
	GetName() string
}

func RunTask(exec Executor, taskName string, tasks map[string]common.Task) error {
	task, ok := tasks[taskName]
	execInfo := fmt.Sprintf("[%s]", exec.GetName())
	taskInfo := fmt.Sprintf("Executing Task: %s", taskName)

	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}

	logger.Info(
		fmt.Sprintf("%s %s", execInfo, taskInfo),
	)

	for _, step := range task.Steps {
		switch step.Type {
		case common.StepTypeCommand:
			str, err := exec.Run(step.Command)
			if err != nil {
				logger.Error(fmt.Sprintf("Command failed: %v\nOutput: %s", err, strings.TrimSpace(str)), exec.GetName())
				return fmt.Errorf("command failed: %v\noutput: %s", err, strings.TrimSpace(str))
			}

			logger.Info(
				fmt.Sprintf("%s -> %s", execInfo, fmt.Sprint(strings.TrimSpace(str))),
			)
		case common.StepTypeTask:
			if err := RunTask(exec, step.Task, tasks); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown step type: %s", step.Type)
		}
	}
	return nil
}

func NewExecutorFromConfig(name string, cfg *config.RemoteHost) (Executor, error) {
	switch cfg.Type {
	case common.ExecutorSSH:
		sshExec, err := newSSHExecutor(name, cfg)
		if err != nil {
			return nil, err
		}
		return sshExec, nil
	case common.ExecutorLocal:
		return newLocalExecutor(name)
	default:
		return nil, fmt.Errorf("unknown executor type: %s", cfg.Type)
	}
}
