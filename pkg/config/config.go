package config

import (
	"fmt"
	"os"

	"github.com/Maru-Yasa/gosong/pkg/executor"
	"gopkg.in/yaml.v3"
)

type RemoteHost struct {
	Type     executor.ExecutorType `yaml:"type"`
	Hostname string                `yaml:"hostname"`
	User     string                `yaml:"user"`
	Port     int                   `yaml:"port,omitempty"`
	KeyPath  string                `yaml:"keyPath,omitempty"`
}

type Config struct {
	Config struct {
		Remote map[string]RemoteHost `yaml:"remote,omitempty"`
	} `yaml:"config"`

	Task map[string]struct {
		Steps []struct {
			Command string `yaml:"command"`
		} `yaml:"steps"`
	}
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
