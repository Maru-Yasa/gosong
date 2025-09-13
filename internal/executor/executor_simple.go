package executor

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
)

// SimpleExecutor is for demo/testing YAML logic only (not production executor)
type SimpleExecutor struct {
	Cwd string
}

func (e *SimpleExecutor) Exec(cmd string) {
	fmt.Printf("[cwd: %s] $ %s\n", e.Cwd, cmd)
}

func RunTaskSimple(cfg *config.Config, exec *SimpleExecutor, taskName string) error {
	task, ok := cfg.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}
	for _, step := range task.Steps {
		switch {
		case step.Cd != "":
			exec.Cwd = step.Cd
			fmt.Printf("Change directory to: %s\n", exec.Cwd)
		case step.Run != "":
			exec.Exec(step.Run)
		case step.Task != "":
			if err := RunTaskSimple(cfg, exec, step.Task); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid step: %+v", step)
		}
	}
	return nil
}
