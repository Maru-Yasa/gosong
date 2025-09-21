package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"
)

func init() {
	RegisterTask(
		"symlink",
		"Update symlink to the latest release",
		func(ctx *Context) error {
			logger.Info("[%s] Updating symlink to latest release...", ctx.Exec.GetName())

			// remove existing symlink
			cmd, err := templateutil.RenderTemplate("rm -f {{.AppPath}}/current", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render rm command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to remove existing symlink: %s", err)
			}

			// create new symlink
			cmd, err = templateutil.RenderTemplate("ln -sfn {{.ReleasePath}} {{.AppPath}}/current", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render ln command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to create symlink: %s", err)
			}

			// verify symlink
			cmd, err = templateutil.RenderTemplate("[[ $(readlink -f {{.AppPath}}/current) == {{.ReleasePath}} ]] || { echo 'symlink error!'; exit 1; }", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render verification command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("symlink verification failed: %s", err)
			}

			// echo success message
			cmd, err = templateutil.RenderTemplate("echo 'symlink current -> {{.ReleasePath}}'", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render echo command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to echo success message: %s", err)
			}

			return nil
		},
	)
}
