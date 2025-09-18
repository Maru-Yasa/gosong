package cmd

import (
	"context"

	"github.com/Maru-Yasa/gosong/actions/daemon"
	internalDeamon "github.com/Maru-Yasa/gosong/internal/daemon"
	"github.com/Maru-Yasa/gosong/pkg/logger"
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
				return internalDeamon.New().Run()
			}

			if err := daemon.Daemon(cli); err != nil {
				logger.NewConsoleLogger().Info("error bro")
				return err
			}
			return nil
		},
	}
}
