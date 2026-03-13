package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/proxy"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/urfave/cli/v3"
)

func DeployInit(cli *cli.Command) error {
	configFilePath := cli.String("config")
	hostName := cli.String("remote")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	logger.Info("[%s] initializing deployment directories with config: %s", "client", configFilePath)

	if hostName == "" {
		for name, remote := range cfg.Config.Remote {
			if err := runInit(name, &remote, &cfg.Config); err != nil {
				logger.Error("[%s] error: %v", name, err)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runInit(hostName, &remote, &cfg.Config)
}

func runInit(hostName string, remote *config.RemoteHost, cfg *config.ConfigRoot) error {
	exec, err := executor.NewExecutorFromConfig(hostName, remote)
	if err != nil {
		return err
	}

	// create deployment directory structure
	dirs := fmt.Sprintf("%s/releases %s/shared /var/log/gosong /var/lib/gosong/apps", cfg.AppPath, cfg.AppPath)
	if _, err := exec.RunRaw("mkdir -p " + dirs); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	logger.Info("[%s] created deployment directories", hostName)

	// if proxy is configured: validate it and create snippet placeholder
	if cfg.Proxy.Type != "" {
		backend, err := proxy.New(cfg.Proxy.Type)
		if err != nil {
			return err
		}
		if err := backend.Validate(exec, cfg.Proxy); err != nil {
			return fmt.Errorf("[%s] proxy validation failed: %w", hostName, err)
		}
		// create empty snippet placeholder so the Caddyfile import doesn't fail on first deploy
		if cfg.Proxy.SnippetPath != "" {
			cmd := fmt.Sprintf("test -f %s || echo '' > %s", cfg.Proxy.SnippetPath, cfg.Proxy.SnippetPath)
			if _, err := exec.RunRaw(cmd); err != nil {
				logger.Warn("[%s] could not create proxy snippet placeholder: %v", hostName, err)
			} else {
				logger.Info("[%s] proxy snippet placeholder ready at %s", hostName, cfg.Proxy.SnippetPath)
			}
		}
	}

	return nil
}
