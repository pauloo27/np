package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TmuxConfig struct {
	WindowCount int `toml:"window_count"`
}

type Config struct {
	Projects map[string]string `toml:"projects"`
	Tmux     TmuxConfig        `toml:"tmux"`
}

func loadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, nixDevProfilesDir, "config.toml")

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
