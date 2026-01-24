package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ProfilesPath string `toml:"profiles_path"`
	Tmux         struct {
		WindowCount int `toml:"window_count"`
	} `toml:"tmux"`
}

func getConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "np", "config.toml")
}

func loadConfig() (*Config, error) {
	configPath := getConfigPath()
	if configPath == "" {
		return nil, os.ErrNotExist
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
