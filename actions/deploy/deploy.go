package deploy

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/proxy"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/templateutil"
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
				logger.Error("[%s] host failed: %v", name, err)
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

func runTaskOnHost(name string, cfg *config.ConfigRoot, remote *config.RemoteHost, taskName string, uTasks map[string]tasks.Task) error {
	exec, err := executor.NewExecutorFromConfig(name, remote)
	if err != nil {
		logger.Error("[%s] Failed to create executor: %v", remote.Hostname, err)
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

	// base overrides
	cfgMap["ReleaseID"] = releaseID
	cfgMap["ReleasePath"] = filepath.Join(cfg.AppPath, "releases", strconv.Itoa(releaseID))

	// inject blue-green port variables when ports are configured
	if len(cfg.App.Ports) >= 2 {
		activePort, err := executor.GetActivePort(exec, cfg.AppPath, cfg.App.Name, cfg.App.Ports)
		if err != nil {
			return fmt.Errorf("failed to determine active port: %w", err)
		}
		greenPort := executor.GetGreenPort(cfg.App.Ports, activePort)
		cfgMap["ActivePort"] = activePort
		cfgMap["GreenPort"] = greenPort
	}

	// inject proxy-switch task dynamically based on proxy type
	if cfg.Proxy.Type != "" {
		backend, err := proxy.New(cfg.Proxy.Type)
		if err != nil {
			return err
		}
		proxyCfg := cfg.Proxy
		if proxyCfg.SnippetPath == "" {
			proxyCfg.SnippetPath = fmt.Sprintf("%s/state/%s.caddy.conf", cfg.AppPath, cfg.App.Name)
		}
		// validate proxy setup before running any tasks
		if err := backend.Validate(exec, proxyCfg); err != nil {
			return fmt.Errorf("[%s] proxy validation failed: %w", name, err)
		}
		greenPort, _ := cfgMap["GreenPort"].(int)
		if uTasks == nil {
			uTasks = make(map[string]tasks.Task)
		}
		uTasks["proxy-switch"] = tasks.Task{
			Description: fmt.Sprintf("Switch %s proxy to port %d", cfg.Proxy.Type, greenPort),
			Steps:       backend.SwitchSteps(greenPort, proxyCfg),
		}
	}

	params := executor.RunTaskParams{
		Exec:     exec,
		Cfg:      cfg,
		CfgMap:   cfgMap,
		TaskName: taskName,
		UTask:    uTasks,
		Cwd:      "",
	}

	if err := executor.RunTask(params); err != nil {
		logger.Error("[%s] Task failed: %v", name, err)
		return err
	}
	return nil
}
