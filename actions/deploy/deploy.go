package deploy

import (
	"fmt"
	"log"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/urfave/cli/v3"
)

func Run(cli *cli.Command) error {
	fmt.Println("Running deploy with config:")

	configFilePath := cli.String("config")
	hostName := cli.String("host")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runOnHost(&remote, "neofetch"); err != nil {
				fmt.Printf("host %s failed: %v\n", name, err)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runOnHost(&remote, "neofetch")
}

func runOnHost(remote *config.RemoteHost, command string) error {
	exec, err := executor.NewExecutorFromConfig(remote)

	if err != nil {
		log.Panic(err)
	}

	output, err := exec.Run(command)

	if err != nil {
		log.Panic(err)
	}

	log.Print(output)

	return nil
}
