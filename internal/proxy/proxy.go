package proxy

import (
	"fmt"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/tasks"
)

// Backend is the interface for proxy server integrations.
// Implement this to add support for nginx, apache, etc.
type Backend interface {
	// Validate checks that the proxy is installed and configured correctly on the host.
	// Returns an error if a required condition is not met.
	Validate(exec executor.Executor, cfg config.ProxyConfig) error

	// SwitchSteps returns the shell command steps to switch traffic to greenPort.
	SwitchSteps(greenPort int, cfg config.ProxyConfig) []tasks.Step
}

// New returns the Backend implementation for the given proxy type.
func New(proxyType string) (Backend, error) {
	switch proxyType {
	case "caddy":
		return &CaddyBackend{}, nil
	default:
		return nil, fmt.Errorf("unsupported proxy type %q (supported: caddy)", proxyType)
	}
}
