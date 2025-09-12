package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/urfave/cli/v3"
)

func Run(cli *cli.Command) error {
	logger.Info(fmt.Sprint("Running deploy with config:"), "client")

	configFilePath := cli.String("config")
	hostName := cli.String("host")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	task := cfg.Task["deploy"]

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runOnHost(name, &remote, &task); err != nil {
				logger.Error(fmt.Sprint("host failed: ", err), name)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runOnHost(hostName, &remote, &task)
}

func runOnHost(name string, remote *config.RemoteHost, task *common.Task) error {
	exec, err := executor.NewExecutorFromConfig(name, remote)

	if err != nil {
		logger.Error(fmt.Sprint("Failed to create executor: ", err), remote.Hostname)
		panic(err)
	}

	exec.RunTask(task)

	return nil
}
