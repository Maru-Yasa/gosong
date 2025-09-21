package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	resource "github.com/Maru-Yasa/gosong/resources"
	"github.com/urfave/cli/v3"
)

func DeployInit(cli *cli.Command) error {
	configFilePath := cli.String("config")
	hostName := cli.String("remote")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	logger.Info("[%s] initialize agent with config: %s", "client", configFilePath)

	if hostName == "" {
		for name, remote := range cfg.Config.Remote {
			err := runInit(name, &remote)
			if err != nil {
				logger.Error("[%s] error: ", name, err)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]

	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runInit(hostName, &remote)
}

func runInit(hostName string, remote *config.RemoteHost) error {
	exec, err := executor.NewExecutorFromConfig(hostName, remote)
	if err != nil {
		return err
	}

	// copy the systemd service config
	cmd := fmt.Sprintf("cat > /tmp/gosong.service <<'EOF'\n%s\nEOF", resource.SystemdConfig)

	if _, err := exec.RunRaw(cmd); err != nil {
		return err
	}

	logger.Info("[%s] copy systemd unit config...", hostName)

	// download the binary thriugh github releases
	fmt.Printf("%s", resource.InitScript)
	result, err := exec.RunRaw(resource.InitScript)
	logger.Info("[%s] download gosong binary: %s", hostName, result)

	return err
}
