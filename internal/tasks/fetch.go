package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"
)

func init() {
	RegisterTask(
		"fetch",
		"Fetch the application source",
		func(ctx *Context) error {
			logger.Info("[%s] Fetching application source...", ctx.Exec.GetName())

			cmd, err := templateutil.RenderTemplate("echo 'Fetching from {{.Source.Type}} -> {{.Source.Url}}'", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render echo command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to echo fetching message: %s", err)
			}

			// create release directory
			cmd, err = templateutil.RenderTemplate("mkdir -p {{.ReleasePath}}", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render mkdir command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to create release directory: %s", err)
			}

			// clone repository
			cmd, err = templateutil.RenderTemplate("git clone --progress --verbose -b {{.Source.Branch}} --depth 1 {{.Source.Url}} {{.ReleasePath}}", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render git clone command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to clone repository: %s", err)
			}

			// remove .git directory
			cmd, err = templateutil.RenderTemplate("rm -rf {{.ReleasePath}}/.git", ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render rm command: %s", err)
			}
			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to remove .git directory: %s", err)
			}

			return nil
		},
	)
}
