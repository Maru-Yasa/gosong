package daemon

import (
	"fmt"
	"log"
	"os"

	"github.com/Maru-Yasa/gosong/internal/agent"
	libdaemon "github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v3"
)

func Daemon(cli *cli.Command) error {
	dContext := libdaemon.Context{
		PidFileName: "gosong.pid",
		PidFilePerm: 0644,
		LogFileName: "daemon-gosong.log",
		LogFilePerm: 0640,
		WorkDir:     "/var/run/gosong",
		Umask:       027,
	}

	d, err := dContext.Reborn()
	if err != nil {
		fmt.Println("Failed to start daemon:", err)
		os.Exit(1)
	}
	if d != nil {
		return nil
	}

	defer dContext.Release()

	dAgent := agent.New()
	if err := dAgent.Run(); err != nil {
		return err
	}

	// serve signal
	err = libdaemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")

	return nil
}
