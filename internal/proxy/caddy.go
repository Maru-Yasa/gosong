package proxy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"github.com/Maru-Yasa/gosong/pkg/logger"
)

// CaddyBackend implements the Backend interface for Caddy web server.
// It manages a Caddyfile snippet per app that the user imports in their main Caddyfile.
type CaddyBackend struct{}

// Validate checks that Caddy is installed and configured correctly on the target host.
func (c *CaddyBackend) Validate(exec executor.Executor, cfg config.ProxyConfig) error {
	// 1. caddy binary must be in PATH
	if _, err := exec.RunRaw("command -v caddy"); err != nil {
		return fmt.Errorf("caddy is not installed or not in PATH on this host")
	}

	// 2. main Caddyfile must exist
	if cfg.CaddyConfig == "" {
		return fmt.Errorf("proxy.caddy_config is not set in your gosong config")
	}
	if _, err := exec.RunRaw(fmt.Sprintf("test -f %s", cfg.CaddyConfig)); err != nil {
		return fmt.Errorf("Caddyfile not found at %s", cfg.CaddyConfig)
	}

	// 3. snippet path must be set (auto-defaulted by deploy before reaching here)
	if cfg.SnippetPath == "" {
		return fmt.Errorf("proxy.snippet_path is not set in your gosong config")
	}

	// 4. ensure snippet directory exists (create it if needed)
	snippetDir := filepath.Dir(cfg.SnippetPath)
	if _, err := exec.RunRaw(fmt.Sprintf("mkdir -p %s", snippetDir)); err != nil {
		return fmt.Errorf("could not create snippet directory %s: %s", snippetDir, err)
	}

	// 5. check that the main Caddyfile imports the snippet (warn only, not fatal)
	out, err := exec.RunRaw(fmt.Sprintf("grep -q 'import %s' %s && echo yes || echo no", cfg.SnippetPath, cfg.CaddyConfig))
	if err == nil && strings.TrimSpace(out) == "no" {
		logger.Warn("Caddyfile at %s does not import %s — add `import %s` to enable proxy switching",
			cfg.CaddyConfig, cfg.SnippetPath, cfg.SnippetPath)
	}

	// 6. validate Caddy config syntax
	if _, err := exec.RunRaw(fmt.Sprintf("caddy validate --config %s", cfg.CaddyConfig)); err != nil {
		return fmt.Errorf("caddy config validation failed for %s: run `caddy validate --config %s` on the host for details",
			cfg.CaddyConfig, cfg.CaddyConfig)
	}

	logger.Info("caddy proxy validation passed")
	return nil
}

// SwitchSteps writes the app's Caddyfile snippet with the new upstream port,
// then reloads Caddy. The user must add `import <snippet_path>` to their main Caddyfile.
func (c *CaddyBackend) SwitchSteps(greenPort int, cfg config.ProxyConfig) []tasks.Step {
	return []tasks.Step{
		{Run: fmt.Sprintf("echo '%s {' > %s", cfg.SiteAddress, cfg.SnippetPath)},
		{Run: fmt.Sprintf("echo '    reverse_proxy localhost:%d' >> %s", greenPort, cfg.SnippetPath)},
		{Run: fmt.Sprintf("echo '}' >> %s", cfg.SnippetPath)},
		{Run: fmt.Sprintf("caddy reload --config %s", cfg.CaddyConfig)},
		{Run: fmt.Sprintf("echo 'Caddy proxy switched to port %d'", greenPort)},
	}
}
