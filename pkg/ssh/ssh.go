package ssh

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	client *ssh.Client
	config *ssh.ClientConfig
	host   string
}
// Run executes a command on the remote SSH server and returns its combined output (stdout + stderr)
func (c *Client) Run(cmd string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	return string(output), err
}

type SSHConfig struct {
	User    string
	Host    string
	Port    uint8
	Keypath string
}

func New(config SSHConfig) (*Client, error) {

	key, err := os.ReadFile(config.Keypath)

	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		return nil, errors.New("failed to parse private key")
	}

	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	host := fmt.Sprintf("%s:%d", config.Host, config.Port)
	fmt.Printf("SSH: Connecting to %s as user %s\n", host, config.User)

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH server: %w", err)
	}

	return &Client{
		client: client,
		config: sshConfig,
		host:   host,
	}, nil
}
