package executor

import (
	"fmt"
	"log"
	"os"

	"github.com/Maru-Yasa/gosong/pkg/config"
	"golang.org/x/crypto/ssh"
)

type SSHExecutor struct {
	Client   *ssh.Client
	Hostname string
}

func newSSHExecutor(cfg *config.RemoteHost) (*SSHExecutor, error) {
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
	log.Printf("SSH: Connecting to %s as user %s\n", host, cfg.User)

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	return &SSHExecutor{
		Client:   client,
		Hostname: cfg.Hostname,
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
