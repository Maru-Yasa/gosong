package deploy

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"github.com/urfave/cli/v3"
)

func Deploy(cli *cli.Command) error {
	logger.Info(
		fmt.Sprintf("[%s] Running deploy with config:", "client"),
	)

	configFilePath := cli.String("config")
	hostName := cli.String("remote")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runTaskOnHost(name, &cfg.Config, &remote, "deploy", cfg.Tasks); err != nil {
				logger.Error(fmt.Sprint("host failed: ", err), name)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runTaskOnHost(hostName, &cfg.Config, &remote, "deploy", cfg.Tasks)
}

func runTaskOnHost(name string, cfg *config.ConfigRoot, remote *config.RemoteHost, taskName string, tasks map[string]tasks.Task) error {
	exec, err := executor.NewExecutorFromConfig(name, remote)
	if err != nil {
		logger.Error(fmt.Sprint("Failed to create executor: ", err), remote.Hostname)
		return err
	}

	lastReleaseID, err := executor.GetLastIDFromHost(exec, cfg.AppPath)

	if err != nil {
		return err
	}

	releaseID := lastReleaseID + 1
	cfgMap, err := templateutil.ToMap(cfg)

	if err != nil {
		return err
	}

	// override config
	cfgMap["ReleaseID"] = releaseID
	cfgMap["ReleasePath"] = filepath.Join(cfg.AppPath, "releases", strconv.Itoa(releaseID))

	params := executor.RunTaskParams{
		Exec:     exec,
		Cfg:      cfg,
		CfgMap:   cfgMap,
		TaskName: taskName,
		UTask:    tasks,
		Cwd:      "",
	}

	if err := executor.RunTask(params); err != nil {
		logger.Error(fmt.Sprint("Task failed: ", err), name)
		return err
	}
	return nil
}
