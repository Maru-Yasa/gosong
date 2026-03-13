package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "gosong-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidMinimal(t *testing.T) {
	path := writeTemp(t, `
config:
  app_path: /opt/myapp
  source:
    type: git
    url: https://github.com/example/app
    branch: main
  remote:
    prod:
      type: ssh
      hostname: example.com
      user: deploy
      port: 22
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Config.AppPath != "/opt/myapp" {
		t.Errorf("AppPath: got %q, want %q", cfg.Config.AppPath, "/opt/myapp")
	}
	if cfg.Config.Source.Type != "git" {
		t.Errorf("Source.Type: got %q, want %q", cfg.Config.Source.Type, "git")
	}
	if cfg.Config.Remote["prod"].Hostname != "example.com" {
		t.Errorf("Remote.prod.Hostname: got %q", cfg.Config.Remote["prod"].Hostname)
	}
}

func TestLoad_WithLocalRemote(t *testing.T) {
	path := writeTemp(t, `
config:
  app_path: /opt/app
  source:
    type: git
    url: https://github.com/example/app
    branch: main
  remote:
    staging:
      type: local
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Config.Remote["staging"].Type != "local" {
		t.Errorf("expected type 'local', got %q", cfg.Config.Remote["staging"].Type)
	}
}

func TestLoad_WithProxyConfig(t *testing.T) {
	path := writeTemp(t, `
config:
  app_path: /opt/myapp
  app:
    name: myapp
    start: "node server.js"
    ports: [3000, 3001]
  source:
    type: git
    url: https://github.com/example/app
    branch: main
  proxy:
    type: caddy
    snippet_path: /etc/caddy/gosong-myapp.conf
    caddy_config: /etc/caddy/Caddyfile
    site_address: myapp.example.com
  remote:
    prod:
      type: ssh
      hostname: example.com
      user: deploy
      port: 22
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Config.Proxy.Type != "caddy" {
		t.Errorf("Proxy.Type: got %q, want %q", cfg.Config.Proxy.Type, "caddy")
	}
	if cfg.Config.Proxy.SnippetPath != "/etc/caddy/gosong-myapp.conf" {
		t.Errorf("Proxy.SnippetPath: got %q", cfg.Config.Proxy.SnippetPath)
	}
	if len(cfg.Config.App.Ports) != 2 {
		t.Errorf("App.Ports: got %v, want [3000, 3001]", cfg.Config.App.Ports)
	}
}

func TestLoad_WithTasks(t *testing.T) {
	path := writeTemp(t, `
config:
  app_path: /opt/app
  source:
    type: git
    url: https://github.com/example/app
    branch: main
  remote:
    prod:
      type: local

tasks:
  deploy:
    description: "Deploy"
    steps:
      - task: lock
      - task: fetch
      - task: unlock
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deploy, ok := cfg.Tasks["deploy"]
	if !ok {
		t.Fatal("expected 'deploy' task")
	}
	if len(deploy.Steps) != 3 {
		t.Errorf("deploy steps: got %d, want 3", len(deploy.Steps))
	}
	if deploy.Steps[0].Task != "lock" {
		t.Errorf("step[0].Task: got %q, want %q", deploy.Steps[0].Task, "lock")
	}
}

func TestLoad_InvalidExecutorType(t *testing.T) {
	path := writeTemp(t, `
config:
  app_path: /opt/app
  source:
    type: git
    url: https://github.com/example/app
    branch: main
  remote:
    prod:
      type: ftp
      hostname: example.com
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid executor type 'ftp'")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := writeTemp(t, `this: is: not: valid: yaml: [`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
