package cmd

import (
	"context"
	"fmt"

	"github.com/Maru-Yasa/gosong/actions/deploy"
	"github.com/urfave/cli/v3"
)

func DeployListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all released deploy",
		Action: func(ctx context.Context, cli *cli.Command) error {
			if err := deploy.DeployList(cli); err != nil {
				return fmt.Errorf("list deployed failed: %w", err)
			}
			return nil
		},
	}
}
