package proxy

import (
	"errors"
	"strings"
	"testing"

	"github.com/Maru-Yasa/gosong/internal/config"
)

// mockExecutor for proxy tests.
type mockExecutor struct {
	name  string
	calls []string
	fn    func(cmd string) (string, error)
}

func (m *mockExecutor) GetName() string { return m.name }
func (m *mockExecutor) RunRaw(cmd string) (string, error) {
	m.calls = append(m.calls, cmd)
	if m.fn != nil {
		return m.fn(cmd)
	}
	return "", nil
}
func (m *mockExecutor) Run(cmd, cwd string) (string, error) { return m.RunRaw(cmd) }

// --- New() factory ---

func TestNew_Caddy(t *testing.T) {
	b, err := New("caddy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := b.(*CaddyBackend); !ok {
		t.Errorf("expected *CaddyBackend, got %T", b)
	}
}

func TestNew_UnknownType(t *testing.T) {
	_, err := New("nginx")
	if err == nil {
		t.Fatal("expected error for unsupported proxy type")
	}
}

func TestNew_EmptyType(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty proxy type")
	}
}

// --- CaddyBackend.SwitchSteps ---

func TestCaddyBackend_SwitchSteps_Count(t *testing.T) {
	b := &CaddyBackend{}
	cfg := config.ProxyConfig{
		SnippetPath: "/etc/caddy/gosong-myapp.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "myapp.example.com",
	}
	steps := b.SwitchSteps(3001, cfg)
	if len(steps) != 5 {
		t.Errorf("expected 5 steps, got %d", len(steps))
	}
}

func TestCaddyBackend_SwitchSteps_ContainsPort(t *testing.T) {
	b := &CaddyBackend{}
	cfg := config.ProxyConfig{
		SnippetPath: "/etc/caddy/app.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "example.com",
	}
	steps := b.SwitchSteps(3001, cfg)

	found := false
	for _, s := range steps {
		if strings.Contains(s.Run, "3001") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected port 3001 to appear in at least one step")
	}
}

func TestCaddyBackend_SwitchSteps_ContainsSiteAddress(t *testing.T) {
	b := &CaddyBackend{}
	cfg := config.ProxyConfig{
		SnippetPath: "/etc/caddy/app.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "mysite.example.com",
	}
	steps := b.SwitchSteps(3000, cfg)

	found := false
	for _, s := range steps {
		if strings.Contains(s.Run, "mysite.example.com") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected site address to appear in at least one step")
	}
}

func TestCaddyBackend_SwitchSteps_ContainsCaddyReload(t *testing.T) {
	b := &CaddyBackend{}
	cfg := config.ProxyConfig{
		SnippetPath: "/etc/caddy/app.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "example.com",
	}
	steps := b.SwitchSteps(3000, cfg)

	found := false
	for _, s := range steps {
		if strings.Contains(s.Run, "caddy reload") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'caddy reload' step")
	}
}

func TestCaddyBackend_SwitchSteps_WritesSnippetPath(t *testing.T) {
	b := &CaddyBackend{}
	cfg := config.ProxyConfig{
		SnippetPath: "/etc/caddy/gosong-myapp.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "example.com",
	}
	steps := b.SwitchSteps(3000, cfg)

	found := false
	for _, s := range steps {
		if strings.Contains(s.Run, "/etc/caddy/gosong-myapp.conf") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected snippet path to appear in steps")
	}
}

// --- CaddyBackend.Validate ---

// allPassMock returns success for all validation commands
func allPassMock() *mockExecutor {
	return &mockExecutor{
		fn: func(cmd string) (string, error) {
			if strings.Contains(cmd, "grep -q") {
				return "yes\n", nil
			}
			return "", nil
		},
	}
}

func validCaddyCfg() config.ProxyConfig {
	return config.ProxyConfig{
		Type:        "caddy",
		SnippetPath: "/etc/caddy/gosong-myapp.conf",
		CaddyConfig: "/etc/caddy/Caddyfile",
		SiteAddress: "example.com",
	}
}

func TestCaddyBackend_Validate_AllPass(t *testing.T) {
	b := &CaddyBackend{}
	exec := allPassMock()
	if err := b.Validate(exec, validCaddyCfg()); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestCaddyBackend_Validate_CaddyNotInPath(t *testing.T) {
	b := &CaddyBackend{}
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		if strings.Contains(cmd, "command -v caddy") {
			return "", errors.New("not found")
		}
		return "", nil
	}}
	err := b.Validate(exec, validCaddyCfg())
	if err == nil {
		t.Fatal("expected error when caddy not in PATH")
	}
	if !strings.Contains(err.Error(), "not installed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCaddyBackend_Validate_MissingCaddyConfig(t *testing.T) {
	b := &CaddyBackend{}
	cfg := validCaddyCfg()
	cfg.CaddyConfig = ""
	err := b.Validate(allPassMock(), cfg)
	if err == nil {
		t.Fatal("expected error for empty caddy_config")
	}
}

func TestCaddyBackend_Validate_CaddyfileNotFound(t *testing.T) {
	b := &CaddyBackend{}
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		if strings.Contains(cmd, "test -f") {
			return "", errors.New("no such file")
		}
		return "", nil
	}}
	err := b.Validate(exec, validCaddyCfg())
	if err == nil {
		t.Fatal("expected error when Caddyfile not found")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCaddyBackend_Validate_MissingSnippetPath(t *testing.T) {
	b := &CaddyBackend{}
	cfg := validCaddyCfg()
	cfg.SnippetPath = ""
	err := b.Validate(allPassMock(), cfg)
	if err == nil {
		t.Fatal("expected error for empty snippet_path")
	}
}

func TestCaddyBackend_Validate_SnippetDirCannotBeCreated(t *testing.T) {
	b := &CaddyBackend{}
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		if strings.Contains(cmd, "mkdir -p") && strings.Contains(cmd, "caddy") {
			return "", errors.New("permission denied")
		}
		return "", nil
	}}
	err := b.Validate(exec, validCaddyCfg())
	if err == nil {
		t.Fatal("expected error when snippet dir cannot be created")
	}
}

func TestCaddyBackend_Validate_CaddyValidateFails(t *testing.T) {
	b := &CaddyBackend{}
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		if strings.Contains(cmd, "caddy validate") {
			return "", errors.New("config parse error")
		}
		if strings.Contains(cmd, "grep -q") {
			return "yes\n", nil
		}
		return "", nil
	}}
	err := b.Validate(exec, validCaddyCfg())
	if err == nil {
		t.Fatal("expected error when caddy validate fails")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCaddyBackend_Validate_MissingImport_IsWarningOnly(t *testing.T) {
	// missing import line in Caddyfile should warn but NOT return an error
	b := &CaddyBackend{}
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		if strings.Contains(cmd, "grep -q") {
			return "no\n", nil // import not found
		}
		return "", nil
	}}
	err := b.Validate(exec, validCaddyCfg())
	if err != nil {
		t.Errorf("missing import should be a warning, not an error; got: %v", err)
	}
}
