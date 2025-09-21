package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
)

var lockFileName = ".gosong.lock"

func init() {
	RegisterTask(
		"deploy:lock",
		"Create a lock file to prevent concurrent deployments",
		func(ctx *Context) error {
			logger.Info("[%s] Locking deployment...", ctx.Exec.GetName())

			// check if already locked
			err := checkLock(ctx, ctx.CfgMap["AppPath"].(string)+"/"+lockFileName)
			if err != nil {
				return fmt.Errorf("deployment is already locked")
			}

			lockFile := ctx.CfgMap["AppPath"].(string) + "/" + lockFileName

			// create the lock file
			err = ctx.Exec.Run("touch "+lockFile, ctx.Cwd)
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
