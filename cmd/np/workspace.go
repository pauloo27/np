package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Workspace struct {
	Projects map[string]string `toml:"projects"`
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
	return filepath.Join(stateDir, "np", "workspace.toml")
}

func loadWorkspace() (*Workspace, error) {
	workspacePath := getWorkspacePath()
	if workspacePath == "" {
		return nil, os.ErrNotExist
	}

	var workspace Workspace
	if _, err := toml.DecodeFile(workspacePath, &workspace); err != nil {
		return nil, err
	}

	return &workspace, nil
}

func saveWorkspace(workspace *Workspace) error {
	workspacePath := getWorkspacePath()
	if workspacePath == "" {
		return os.ErrInvalid
	}

	workspaceDir := filepath.Dir(workspacePath)

	if err := os.MkdirAll(workspaceDir, 0750); err != nil {
		return err
	}

	file, err := os.Create(workspacePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(workspace)
}
