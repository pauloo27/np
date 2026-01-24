package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Workspace struct {
	Projects map[string]string `toml:"projects"`
}

func loadWorkspace() (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	workspacePath := filepath.Join(homeDir, ".local", "state", "np", "workspace.toml")

	var workspace Workspace
	if _, err := toml.DecodeFile(workspacePath, &workspace); err != nil {
		return nil, err
	}

	return &workspace, nil
}

func saveWorkspace(workspace *Workspace) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	workspacePath := filepath.Join(homeDir, ".local", "state", "np", "workspace.toml")
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
