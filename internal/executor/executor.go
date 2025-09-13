package executor

import (
	"fmt"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

type Executor interface {
	Run(cmd string, cwd string) (string, error)
	GetName() string
}

// RunTask runs a task by name, supporting step fields: cd, run, task. Maintains cwd state across tasks.
func RunTask(exec Executor, taskName string, tasks map[string]common.Task, cwd string) error {
	task, ok := tasks[taskName]
	execInfo := fmt.Sprintf("[%s]", exec.GetName())
	taskInfo := fmt.Sprintf("Executing Task: %s", taskName)

	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}

	logger.Info(
		fmt.Sprintf("%s %s", execInfo, taskInfo),
	)

	currentCwd := cwd
	for _, step := range task.Steps {
		switch {
		case step.Cd != "":
			currentCwd = step.Cd
			logger.Info(fmt.Sprintf("Change directory to: %s", currentCwd))
		case step.Run != "":
			str, err := exec.Run(step.Run, currentCwd)
			if err != nil {
				logger.Error(fmt.Sprintf("Command failed: %v\nOutput: %s", err, strings.TrimSpace(str)), exec.GetName())
				return fmt.Errorf("command failed: %v\noutput: %s", err, strings.TrimSpace(str))
			}
			logger.Info(fmt.Sprint(strings.TrimSpace(str)))
		case step.Task != "":
			if err := RunTask(exec, step.Task, tasks, currentCwd); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid step: %+v", step)
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
