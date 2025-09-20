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
				Name:    "remote",
				Aliases: []string{"r"},
				Usage:   "Remote Host",
			},
		},
		Commands: []*cli.Command{
			DeployListCommand(),
			DeployInitCommand(),
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			if err := deploy.Deploy(cli); err != nil {
				return fmt.Errorf("deploy failed: %w", err)
			}
			return nil
		},
	}
}

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

func DeployInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize deployment on remote host(s)",
		Action: func(ctx context.Context, cli *cli.Command) error {
			if err := deploy.DeployInit(cli); err != nil {
				return fmt.Errorf("deploy init failed: %w", err)
			}
			return nil
		},
	}
}
