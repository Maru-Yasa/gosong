package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/urfave/cli/v3"
)

func DeployList(cli *cli.Command) error {

	configFilePath := cli.String("config")
	hostName := cli.String("remote")

	cfg, err := config.Load(configFilePath)

	if err != nil {
		return err
	}

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runOnHost(name, &cfg.Config, &remote, "show_releases", cfg.Tasks); err != nil {
				logger.Error(fmt.Sprint("host failed: ", err), name)
			}
		}
		return nil
	}

	remote := cfg.Config.Remote[hostName]

	return runOnHost(hostName, &cfg.Config, &remote, "show_releases", cfg.Tasks)
}
