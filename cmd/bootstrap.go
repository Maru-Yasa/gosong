package cmd

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func Execute() {
	app := &cli.Command{
		Name:  "gosong",
		Usage: "zero-downtime deploy tool",
		Commands: []*cli.Command{
			DeployCommand(),
			DaemonCommand(),
			ProcessCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
