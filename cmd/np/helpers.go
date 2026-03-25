package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

func getProfilesPath() string {
	if cfg.ProfilesPath != "" {
		return cfg.ProfilesPath
	}

	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "nix-conf/dev")
}

func profileCompletion(includeLocal bool) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profilesPath := getProfilesPath()
		profiles, err := getAvailableProfiles(profilesPath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if includeLocal {
			profiles = append([]string{"none", "local"}, profiles...)
		}
		return profiles, cobra.ShellCompDirectiveNoFileComp
	}
}

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

	if project, exists := workspace.Projects[cwd]; exists {
		return project.Profile, true
	}

	// TODO: recursive check?

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

func getBinPath(name string) (string, error) {
	binPath, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s binary not found in PATH: %w", name, err)
	}
	return binPath, nil
}

func getShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	return shell
}

func resolveVariation() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	if project, exists := workspace.Projects[cwd]; exists {
		return project.Variation
	}
	return ""
}

func resolveProfile(args []string, nixDevProfilesPath string) (profile string, useLocalFlake bool, ok bool) {
	if len(args) > 0 {
		profile = args[0]
		if profile == "none" {
			return profile, false, true
		}
		if profile == "local" {
			return profile, true, true
		}
		return profile, false, true
	}

	if envProfile := os.Getenv("USING_NIX_DEV"); envProfile != "" {
		profile = envProfile
		if profile == "none" {
			return profile, false, true
		}
		if profile == "local" {
			return profile, true, true
		}
		return profile, false, true
	}

	detectedProfile, found := determineProfileName()
	if found {
		profile = detectedProfile
		if profile == "none" {
			return profile, false, true
		}
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

func buildTmuxWindows(windowCount int, windowCommands []string) []*config.TmuxWindow {
	totalCount := windowCount + len(windowCommands)
	if totalCount == 0 {
		windowCount = 1
		totalCount = 1
	}

	windows := make([]*config.TmuxWindow, 0, totalCount)

	for _, cmd := range windowCommands {
		windows = append(windows, &config.TmuxWindow{Command: cmd})
	}

	for range windowCount {
		windows = append(windows, &config.TmuxWindow{})
	}

	return windows
}
