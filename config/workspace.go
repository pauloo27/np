package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Workspace struct {
	Projects map[string]string `yaml:"projects"`
	path     string            `yaml:"-"` // configured workspace path
}

func (c *Config) GetWorkspacePath() string {
	if c != nil && c.WorkspacePath != "" {
		return c.WorkspacePath
	}

	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		stateDir = filepath.Join(homeDir, ".local", "state")
	}
	return filepath.Join(stateDir, "np", "workspace.yaml")
}

func LoadWorkspace(cfg *Config) (*Workspace, error) {
	workspacePath := cfg.GetWorkspacePath()
	if workspacePath == "" {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(workspacePath)
	if err != nil {
		return nil, err
	}

	var workspace Workspace
	if err := yaml.Unmarshal(data, &workspace); err != nil {
		return nil, err
	}

	workspace.path = workspacePath
	return &workspace, nil
}

func (w *Workspace) Save() error {
	if w.path == "" {
		return os.ErrInvalid
	}

	workspaceDir := filepath.Dir(w.path)

	if err := os.MkdirAll(workspaceDir, 0750); err != nil {
		return err
	}

	data, err := yaml.Marshal(w)
	if err != nil {
		return err
	}

	return os.WriteFile(w.path, data, 0640)
}
