package cmd

import (
	"context"
	"fmt"

	"github.com/Maru-Yasa/gosong/actions/deploy"
	"github.com/urfave/cli/v3"
)

func DeployCommand() *cli.Command {
	return &cli.Command{
		Name:  "deploy",
		Usage: "Deploy an application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Config file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"h"},
				Usage:   "Remote Host",
			},
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			if err := deploy.Run(cli); err != nil {
				return fmt.Errorf("deploy failed: %w", err)
			}
			return nil
		},
	}
}
