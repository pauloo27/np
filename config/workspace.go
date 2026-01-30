package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Workspace struct {
	Projects map[string]string `yaml:"projects"`
}

func getWorkspacePath() string {
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

func LoadWorkspace() (*Workspace, error) {
	workspacePath := getWorkspacePath()
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

	return &workspace, nil
}

func SaveWorkspace(workspace *Workspace) error {
	workspacePath := getWorkspacePath()
	if workspacePath == "" {
		return os.ErrInvalid
	}

	workspaceDir := filepath.Dir(workspacePath)

	if err := os.MkdirAll(workspaceDir, 0750); err != nil {
		return err
	}

	data, err := yaml.Marshal(workspace)
	if err != nil {
		return err
	}

	return os.WriteFile(workspacePath, data, 0640)
}
