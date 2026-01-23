package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func listAvailableProfiles(profilesPath string) {
	availableProfiles, err := getAvailableProfiles(profilesPath)
	if err == nil && len(availableProfiles) > 0 {
		fmt.Fprintf(os.Stderr, "available profiles:")
		for _, p := range availableProfiles {
			fmt.Fprintf(os.Stderr, " %s", p)
		}
		fmt.Fprintln(os.Stderr)
	}
}

func determineProfileName() (string, bool) {
	if _, err := os.Stat("flake.nix"); err == nil {
		return "local", true
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	config, err := loadConfig()
	if err != nil {
		return "", false
	}

	if profile, exists := config.Projects[cwd]; exists {
		return profile, true
	}

	absPath, err := filepath.EvalSymlinks(cwd)
	if err == nil {
		if profile, exists := config.Projects[absPath]; exists {
			return profile, true
		}
	}

	return "", false
}

func getAvailableProfiles(profilesPath string) ([]string, error) {
	entries, err := os.ReadDir(profilesPath)
	if err != nil {
		return nil, err
	}

	var profiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			profiles = append(profiles, entry.Name())
		}
	}
	return profiles, nil
}
