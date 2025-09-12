package executor

import (
	"fmt"
	"os"
	"strings"

	"github.com/Maru-Yasa/gosong/pkg/logger"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSHExecutor struct {
	Client *ssh.Client
	Name   string
}

func newSSHExecutor(name string, cfg *config.RemoteHost) (*SSHExecutor, error) {
	key, err := os.ReadFile(cfg.KeyPath)

	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	host := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)
	logger.Info(fmt.Sprint("SSH: Connecting as user ", cfg.User, " (", host, ")"), name)

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	return &SSHExecutor{
		Client: client,
		Name:   name,
	}, nil
}

func (s *SSHExecutor) Run(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

func (s *SSHExecutor) RunTask(task *common.Task) {
	for _, step := range task.Steps {
		str, err := s.Run(step.Command)

		if err != nil {
			logger.Error(fmt.Sprint("SSH command failed: ", err), s.Name)
			panic(err)
		}

		logger.Info(fmt.Sprint("SSH command output: ", strings.TrimSpace(str)), s.Name)
	}
}
