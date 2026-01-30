package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	ProfilesPath  string `yaml:"profiles_path"`
	WorkspacePath string `yaml:"workspace_path,omitempty"`
	Tmux          struct {
		WindowCount int `yaml:"window_count"`
	} `yaml:"tmux"`
}

func GetConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "np", "config.yaml")
}

func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	if configPath == "" {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
