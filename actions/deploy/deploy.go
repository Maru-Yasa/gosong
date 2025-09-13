package deploy

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"

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

	lastReleaseID, err := getLastIDFromHost(exec, cfg.AppPath)

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

func getLastIDFromHost(exec executor.Executor, appPath string) (int, error) {
	script := fmt.Sprintf(`
		if [ ! -d %s/releases ]; then
			echo 0
			exit 0
		fi

		last_id=$(find %s/releases -maxdepth 1 -type d -printf "%%f\n" \
			| grep -E '^[0-9]+$' \
			| sort -n \
			| tail -1)

		if [ -z "$last_id" ]; then
			echo 0
		else
			echo "$last_id"
		fi
	`, appPath, appPath)

	outputFromExec, err := exec.RunRaw(script)
	if err != nil {
		return 0, err
	}

	id, convErr := strconv.Atoi(strings.TrimSpace(outputFromExec))
	if convErr != nil {
		return 0, convErr
	}

	return id, nil
}
