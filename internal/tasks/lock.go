package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
)

var lockFileName = ".gosong.lock"

func init() {
	RegisterTask(
		"lock",
		"Create a lock file to prevent concurrent deployments",
		func(ctx *Context) error {
			logger.Info("[%s] Locking deployment...", ctx.Exec.GetName())

			// check if already locked
			lockFileCmd, err := ctx.RenderCmd("{{.AppPath}}/"+lockFileName, ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render lock file path: %s", err)
			}

			err = checkLock(ctx, lockFileCmd)
			if err != nil {
				return fmt.Errorf("deployment is already locked")
			}

			// create the lock file
			cmd, err := ctx.RenderCmd("touch {{.AppPath}}/"+lockFileName, ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render touch command: %s", err)
			}

			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to create lock file: %s", err)
			}

			return nil
		},
	)
}

func checkLock(ctx *Context, lockFile string) error {
	// check if already locked
	script := fmt.Sprintf("if [ -f %s ]; then exit 1; fi", lockFile)
	err := ctx.Exec.Run(script, ctx.Cwd)

	if err != nil {
		return err
	}

	return nil
}
