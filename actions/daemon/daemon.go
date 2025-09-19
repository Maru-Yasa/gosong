package daemon

import (
	"os"

	"github.com/Maru-Yasa/gosong/internal/daemon"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	libdaemon "github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v3"
)

func Daemon(cli *cli.Command) error {
	dContext := libdaemon.Context{
		PidFileName: "gosong.pid",
		PidFilePerm: 0644,
		LogFileName: "daemon-gosong.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	d, err := dContext.Reborn()
	if err != nil {
		fileLogger := logger.NewConsoleLogger()
		fileLogger.Error("Failed to start daemon: %v", err)
		os.Exit(1)
	}
	if d != nil {
		return nil
	}

	defer dContext.Release()

	dAgent := daemon.New()
	if err := dAgent.Run(); err != nil {
		return err
	}

	// serve signal
	fileLogger := logger.NewFileLogger()
	err = libdaemon.ServeSignals()
	if err != nil {
		fileLogger.Error("Error: %s", err.Error())
		fileLogger.Sync()
	}

	fileLogger.Info("daemon terminated")
	fileLogger.Sync()

	return nil
}
