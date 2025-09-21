package executor

import (
	"fmt"
	"os"
	"time"

	"github.com/Maru-Yasa/gosong/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSHExecutor struct {
	Client *ssh.Client
	Name   string
}

func (s *SSHExecutor) GetName() string {
	return s.Name
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
		Timeout:         5 * time.Second, // set timeout
	}

	host := fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port)

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	return &SSHExecutor{
		Client: client,
		Name:   name,
	}, nil
}

func (s *SSHExecutor) RunRaw(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()
	fullCmd := cmd
	out, err := session.CombinedOutput(fullCmd)

	return string(out), err
}

func (s *SSHExecutor) Run(cmd string, cwd string) error {
	session, err := s.Client.NewSession()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err != nil {
		return err
	}

	defer session.Close()
	fullCmd := cmd
	if cwd != "" {
		fullCmd = "cd " + cwd + " && " + cmd
	}

	err = session.Run(fullCmd)

	return err
}
