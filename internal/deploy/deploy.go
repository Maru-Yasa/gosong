package deploy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/pkg/config"
	"github.com/Maru-Yasa/gosong/pkg/executor"
	"github.com/Maru-Yasa/gosong/pkg/ssh"
	"github.com/urfave/cli/v3"
)

func Run(cli *cli.Command) error {
	fmt.Println("Running deploy with config:")

	configFilePath := cli.String("config")
	hostName := cli.String("host")

	cfg, err := config.Load(configFilePath)
	if err != nil {
		return err
	}

	if hostName == "" {
		// loop all hosts
		for name, remote := range cfg.Config.Remote {
			if err := runOnHost(name, remote, "neofetch"); err != nil {
				fmt.Printf("host %s failed: %v\n", name, err)
			}
		}
		return nil
	}

	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return fmt.Errorf("host %s not found in config", hostName)
	}

	return runOnHost(hostName, remote, "neofetch")
}

func runOnHost(name string, remote config.RemoteHost, command string) error {
	sshConfig := ssh.SSHConfig{
		Host:    remote.Hostname,
		Port:    uint8(remote.Port),
		User:    remote.User,
		Keypath: remote.KeyPath,
	}

	sshClient, err := ssh.New(sshConfig)
	if err != nil {
		return fmt.Errorf("failed to create SSH client for %s: %w", name, err)
	}

	fmt.Printf("executing on host: %s (%s)\n", name, remote.Hostname)
	executor.ExecuteCommand(sshClient, command)
	return nil
}
