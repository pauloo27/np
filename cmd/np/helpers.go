package main

import (
	"fmt"
	"os"
	"os/exec"
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

func runCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func resolveProfile(args []string, nixDevProfilesPath string) (profile string, useLocalFlake bool, ok bool) {
	if len(args) > 0 {
		profile = args[0]
		if profile == "local" {
			return profile, true, true
		}
		return profile, false, true
	}

	if envProfile := os.Getenv("USING_NIX_DEV"); envProfile != "" {
		profile = envProfile
		if profile == "local" {
			return profile, true, true
		}
		return profile, false, true
	}

	detectedProfile, found := determineProfileName()
	if found {
		profile = detectedProfile
		if profile == "local" {
			return profile, true, true
		}
		return profile, false, true
	}

	fmt.Fprintf(os.Stderr, "no profile specified and could not determine one automatically\n")
	fmt.Fprintf(os.Stderr, "use 'np set <profile>' to set a profile for this directory\n")
	listAvailableProfiles(nixDevProfilesPath)
	return "", false, false
}
