package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
)

func init() {
	RegisterTask(
		"unlock",
		"Remove lock file to allow deployments",
		func(ctx *Context) error {
			logger.Info("[%s] Unlocking deployment...", ctx.Exec.GetName())

			// remove the lock file
			cmd, err := ctx.RenderCmd("rm -f {{.AppPath}}/.gosong.lock", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render rm command: %s", err)
			}

			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to remove lock file: %s", err)
			}

			return nil
		},
	)
}
