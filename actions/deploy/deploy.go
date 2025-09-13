package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/logger"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"github.com/urfave/cli/v3"
)

func Run(cli *cli.Command) error {
	logger.Info(
		fmt.Sprintf("[%s] Running deploy with config:", "client"),
	)

	configFilePath := cli.String("config")
	hostName := cli.String("host")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runOnHost(name, &cfg.Config, &remote, "deploy", cfg.Tasks); err != nil {
				logger.Error(fmt.Sprint("host failed: ", err), name)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runOnHost(hostName, &cfg.Config, &remote, "deploy", cfg.Tasks)
}

func runOnHost(name string, cfg *config.ConfigRoot, remote *config.RemoteHost, taskName string, tasks map[string]tasks.Task) error {
	exec, err := executor.NewExecutorFromConfig(name, remote)
	if err != nil {
		logger.Error(fmt.Sprint("Failed to create executor: ", err), remote.Hostname)
		return err
	}
	if err := executor.RunTask(exec, cfg, taskName, tasks, ""); err != nil {
		logger.Error(fmt.Sprint("Task failed: ", err), name)
		return err
	}
	return nil
}
