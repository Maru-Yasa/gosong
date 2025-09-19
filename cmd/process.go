package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/daemon"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/unixsocket"
	"github.com/urfave/cli/v3"
)

func ProcessCommand() *cli.Command {
	return &cli.Command{
		Name:  "process",
		Usage: "process <command>",
		Commands: []*cli.Command{
			ProcessStartCommand(),
			ProcessStatusCommand(),
			ProcessStopCommand(),
		},
	}
}

func ProcessStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "<app_name> -- [args...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bin",
				Aliases:  []string{"b"},
				Usage:    "Path to binary",
				Required: false,
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Usage:    "Port to run the app on",
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			consoleLogger := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				consoleLogger.Error("need <app_name>")
				return fmt.Errorf("need <app_name>")
			}

			app := c.Args().First()
			bin := c.String("bin")
			port := c.Int("port")
			args := strings.Join(c.Args().Slice()[1:], ",")

			conn := unixsocket.NewClient("/tmp/gosong.sock")

			payload := map[string]string{
				"action": string(daemon.ProcessActionStart),
				"app":    app,
				"bin":    bin,
				"port":   strconv.Itoa(port),
				"args":   args,
			}

			result, err := conn.Send(payload)

			if err != nil {
				return fmt.Errorf("failed to send command to daemon: %w", err)
			}

			consoleLogger.Info(result)

			return nil
		},
	}
}

func ProcessStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "<app_name>",
		Action: func(ctx context.Context, c *cli.Command) error {
			consoleLogger := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				consoleLogger.Error("need <app_name>")
				return fmt.Errorf("need <app_name>")
			}

			conn := unixsocket.NewClient("/tmp/gosong.sock")
			app := c.Args().First()
			payload := map[string]string{
				"action": string(daemon.ProcessActionStatus),
				"app":    app,
			}

			result, err := conn.Send(payload)

			if err != nil {
				return err
			}

			consoleLogger.Info(result)

			return nil
		},
	}
}

func ProcessStopCommand() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "<app_name>",
		Action: func(ctx context.Context, c *cli.Command) error {
			consoleLogger := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				consoleLogger.Error("need <app_name>")
				return fmt.Errorf("need <app_name>")
			}

			conn := unixsocket.NewClient("/tmp/gosong.sock")
			app := c.Args().First()
			payload := map[string]string{
				"action": string(daemon.ProcessActionStop),
				"app":    app,
			}

			result, err := conn.Send(payload)

			if err != nil {
				return err
			}

			consoleLogger.Info(result)

			return nil
		},
	}
}
