package tasks

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"
)

func init() {
	RegisterTask(
		"show_releases",
		"Show releases",
		func(ctx *Context) error {
			logger.Info("[%s] Showing releases...", ctx.Exec.GetName())

			script := `
				current_id=$(readlink -f {{.AppPath}}/current 2>/dev/null | xargs basename)

				ls -1 {{.AppPath}}/releases | grep -E '^[0-9]+
 | sort -n -r | while read id; do
					if [ "$id" = "$current_id" ]; then
						echo "$id (current)"
					else
						echo "$id"
					fi
				done
			`

			cmd, err := templateutil.RenderTemplate(script, ctx.CfgMap)
			if err != nil {
				return fmt.Errorf("failed to render show releases script: %s", err)
			}

			err = ctx.Exec.Run(cmd, ctx.Cwd)
			if err != nil {
				return fmt.Errorf("failed to show releases: %s", err)
			}

			return nil
		},
	)
}
