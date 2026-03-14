package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Maru-Yasa/gosong/internal/config"
	"golang.org/x/crypto/ssh"
)

type SSHExecutor struct {
	Client    *ssh.Client
	Name      string
	ShellInit string // extra init sourced before every command
}

func (s *SSHExecutor) GetName() string {
	return s.Name
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func newSSHExecutor(name string, cfg *config.RemoteHost) (*SSHExecutor, error) {
	key, err := os.ReadFile(expandPath(cfg.KeyPath))
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
		Client:    client,
		Name:      name,
		ShellInit: cfg.ShellInit,
	}, nil
}

// preamble returns shell init commands that run before every remote command.
// It always sources nvm and asdf (guarded so they are safe when absent),
// plus any user-defined shell_init from the config.
func (s *SSHExecutor) preamble() string {
	parts := []string{
		// nvm — most common Node version manager
		`export NVM_DIR="$HOME/.nvm"`,
		`[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"`,
		// asdf — multi-runtime version manager
		`[ -s "$HOME/.asdf/asdf.sh" ] && . "$HOME/.asdf/asdf.sh"`,
	}
	if s.ShellInit != "" {
		parts = append(parts, s.ShellInit)
	}
	return strings.Join(parts, "; ")
}

// loginShell wraps cmd in a bash login shell with the version-manager preamble
// so tools installed via nvm/asdf/etc. are available in PATH.
func (s *SSHExecutor) loginShell(cmd string) string {
	escaped := strings.ReplaceAll(cmd, "'", `'\''`)
	return "bash -l -c '" + s.preamble() + "; " + escaped + "'"
}

func (s *SSHExecutor) RunRaw(cmd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(s.loginShell(cmd))
	return string(out), err
}

func (s *SSHExecutor) Run(cmd string, cwd string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	defer session.Close()

	fullCmd := cmd
	if cwd != "" {
		fullCmd = "cd " + cwd + " && " + cmd
	}

	err = session.Run(s.loginShell(fullCmd))
	return fmt.Sprintf("running command -> %s", fullCmd), err
}
