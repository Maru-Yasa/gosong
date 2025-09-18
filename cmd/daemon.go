package cmd

import (
	"context"
	"fmt"

	"github.com/Maru-Yasa/gosong/actions/daemon"
	"github.com/Maru-Yasa/gosong/internal/agent"
	"github.com/urfave/cli/v3"
)

func DaemonCommand() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Run Daemon",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "foreground",
				Aliases: []string{"f"},
				Usage:   "Foreground",
			},
		},
		Action: func(ctx context.Context, cli *cli.Command) error {
			if cli.Bool("foreground") {
				return agent.New().Run()
			}

			if err := daemon.Daemon(cli); err != nil {
				fmt.Println("error bro")
				return fmt.Errorf("daemon failed to start: %w", err)
			}
			return nil
		},
	}
}
