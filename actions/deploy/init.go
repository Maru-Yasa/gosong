package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
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
	script := `
REPO="Maru-Yasa/gosong"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
esac

ASSET_URL=$(curl -s "https://api.github.com/repos/$REPO/releases" \
  | grep "browser_download_url" \
  | grep "$OS-$ARCH.zip" \
  | head -n 1 \
  | cut -d '"' -f 4)

echo "Downloading: $ASSET_URL"
curl -L "$ASSET_URL" -o gosong-$OS-$ARCH.zip
unzip -o gosong-$OS-$ARCH.zip -d /opt/gosong

mv /opt/gosong/gosong-$OS-$ARCH /opt/gosong/gosong
chmod +x /opt/gosong/gosong
	`
	exec, err := executor.NewExecutorFromConfig(hostName, remote)
	if err != nil {
		return err
	}

	result, err := exec.RunRaw(script)
	logger.Info("[%s] init gosong: %s", hostName, result)
	return err
}
