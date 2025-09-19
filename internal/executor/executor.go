package executor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"
)

// RunTaskParams is used to pass parameters to RunTask
type RunTaskParams struct {
	Exec     Executor
	Cfg      *config.ConfigRoot
	CfgMap   map[string]any
	TaskName string
	UTask    map[string]tasks.Task
	Cwd      string
}

type Executor interface {
	RunRaw(cmd string) (string, error)
	Run(cmd string, cwd string) (string, error)
	GetName() string
}

func RunTask(params RunTaskParams) error {
	// searching for task, whatever on user defined tasks or mine
	task, err := tasks.FindTask(params.TaskName, params.UTask)
	if err != nil {
		return err
	}

	// for logging
	execInfo := fmt.Sprintf("[%s]", params.Exec.GetName())
	taskInfo := fmt.Sprintf("Executing Task: %s", params.TaskName)

	logger.Info(
		fmt.Sprintf("%s %s", execInfo, taskInfo),
	)

	currentCwd := params.Cwd
	for _, step := range task.Steps {
		switch {
		case step.Cd != "":
			cmdCd, err := templateutil.RenderTemplate(step.Cd, params.CfgMap)
			if err != nil {
				return fmt.Errorf("command failed to render: %s", err)
			}
			currentCwd = cmdCd
			logger.Info(fmt.Sprintf("Change directory to: %s", currentCwd))
		case step.Run != "":
			cmdRun, err := templateutil.RenderTemplate(step.Run, params.CfgMap)

			if err != nil {
				return fmt.Errorf("command failed to render: %s", err)
			}
			str, err := params.Exec.Run(cmdRun, currentCwd)
			if err != nil {
				logger.Error(fmt.Sprintf("Command failed: %v\nOutput: %s", err, strings.TrimSpace(str)), params.Exec.GetName())
				return fmt.Errorf("command failed: %v\noutput: %s", err, strings.TrimSpace(str))
			}
			logger.Info(fmt.Sprint(strings.TrimSpace(str)))
		case step.Task != "":
			if err := RunTask(RunTaskParams{
				Exec:     params.Exec,
				Cfg:      params.Cfg,
				CfgMap:   params.CfgMap,
				TaskName: step.Task,
				UTask:    params.UTask,
				Cwd:      currentCwd,
			}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid step: %+v", step)
		}
	}
	return nil
}

func GetLastIDFromHost(exec Executor, appPath string) (int, error) {
	script := fmt.Sprintf(`
		if [ ! -d %s/releases ]; then
			echo 0
			exit 0
		fi

		last_id=$(find %s/releases -maxdepth 1 -type d -printf "%%f\n" \
			| grep -E '^[0-9]+$' \
			| sort -n \
			| tail -1)

		if [ -z "$last_id" ]; then
			echo 0
		else
			echo "$last_id"
		fi
	`, appPath, appPath)

	outputFromExec, err := exec.RunRaw(script)
	if err != nil {
		return 0, err
	}

	id, convErr := strconv.Atoi(strings.TrimSpace(outputFromExec))
	if convErr != nil {
		return 0, convErr
	}

	return id, nil
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
