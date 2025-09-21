package tasks

import (
	"fmt"
	"os"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

type TaskFunc func(ctx *Context) error

type Context struct {
	CfgMap map[string]any
	Exec   executor.Executor
	Cwd    string

	RenderCmd func(cmd string, data map[string]any) (string, error)
}

type BTask struct {
	Name        string
	Description string

	Run func(ctx *Context) error
}

var BuiltInTasks = map[string]BTask{}

func RegisterTask(name string, desc string, fn TaskFunc) {
	BuiltInTasks[name] = BTask{
		Name:        name,
		Description: desc,
		Run:         fn,
	}
}

func FindAndRun(name string, uTasks map[string]common.UTask, ctx *Context) error {
	// handle built-in tasks first
	bTask, ok := BuiltInTasks[name]
	if ok {
		err := bTask.Run(ctx)
		if err != nil {
			logger.Error("[%s] Task failed: %v", ctx.Exec.GetName(), err)
			os.Exit(1)
		}

		return nil
	}

	// then user defined tasks
	uTask, ok := uTasks[name]
	if !ok {
		return fmt.Errorf("task '%s' not found", name)
	}

	// run the steps and call recursively if needed
	for _, step := range uTask.Steps {
		if step.Task != "" {
			err := FindAndRun(step.Task, uTasks, ctx)
			if err != nil {
				return err
			}
		} else if step.Run != "" {
			rCmd, err := ctx.RenderCmd(step.Run, ctx.CfgMap)

			if err != nil {
				return fmt.Errorf("command failed to render: %s", err)
			}

			err = ctx.Exec.Run(rCmd, ctx.Cwd)

			if err != nil {
				return err
			}
		} else if step.Cd != "" {
			cmdCd, err := ctx.RenderCmd(step.Cd, ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("command failed to render: %s", err)
			}
			ctx.Cwd = cmdCd
		}
	}

	return nil
}
