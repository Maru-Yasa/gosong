package config

import (
	"fmt"
	"os"

	"github.com/Maru-Yasa/gosong/internal/common"
	"github.com/Maru-Yasa/gosong/internal/tasks"
	"gopkg.in/yaml.v3"
)

type RemoteHost struct {
	Type     common.ExecutorType `yaml:"type"`
	Hostname string              `yaml:"hostname"`
	User     string              `yaml:"user"`
	Port     int                 `yaml:"port,omitempty"`
	KeyPath  string              `yaml:"keyPath,omitempty"`
}

type Source struct {
	Type   string `yaml:"type"` // git | local
	Url    string `yaml:"url,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

type ConfigRoot struct {
	Remote  map[string]RemoteHost `yaml:"remote,omitempty"`
	AppPath string                `yaml:"app_path"`
	Source  Source                `yaml:"source"`
}

type Config struct {
	Config ConfigRoot            `yaml:"config"`
	Tasks  map[string]tasks.Task `yaml:"tasks"`
}

func Load(filePath string) (*Config, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		return nil, err
	}

	for name, remote := range cfg.Config.Remote {
		if !remote.Type.IsValid() {
			return nil, fmt.Errorf("invalid type %q for host %s", remote.Type, name)
		}
	}

	return &cfg, nil
}
