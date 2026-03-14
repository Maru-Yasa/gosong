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

	// create all directories under app_path (no root required)
	appDirs := fmt.Sprintf("%s/releases %s/shared %s/logs %s/state", cfg.AppPath, cfg.AppPath, cfg.AppPath, cfg.AppPath)
	if out, err := exec.RunRaw("mkdir -p " + appDirs); err != nil {
		return fmt.Errorf("failed to create app directories: %w\n%s", err, out)
	}
	logger.Info("[%s] created deployment directories", hostName)

	// if proxy is configured: validate it and create snippet placeholder
	if cfg.Proxy.Type != "" {
		backend, err := proxy.New(cfg.Proxy.Type)
		if err != nil {
			return err
		}
		proxyCfg := cfg.Proxy
		if proxyCfg.SnippetPath == "" {
			proxyCfg.SnippetPath = fmt.Sprintf("%s/state/%s.caddy.conf", cfg.AppPath, cfg.App.Name)
		}
		if err := backend.Validate(exec, proxyCfg); err != nil {
			return fmt.Errorf("[%s] proxy validation failed: %w", hostName, err)
		}
		// create empty snippet placeholder so the Caddyfile import doesn't fail on first deploy
		cmd := fmt.Sprintf("test -f %s || echo '' > %s", proxyCfg.SnippetPath, proxyCfg.SnippetPath)
		if _, err := exec.RunRaw(cmd); err != nil {
			logger.Warn("[%s] could not create proxy snippet placeholder: %v", hostName, err)
		} else {
			logger.Info("[%s] proxy snippet placeholder ready at %s — add `import %s` to your Caddyfile", hostName, proxyCfg.SnippetPath, proxyCfg.SnippetPath)
		}
	}

	return nil
}
