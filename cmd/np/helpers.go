package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getProfilesPath() string {
	if config.ProfilesPath != "" {
		return config.ProfilesPath
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
			profiles = append([]string{"local"}, profiles...)
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

	if profile, exists := workspace.Projects[cwd]; exists {
		return profile, true
	}

	absPath, err := filepath.EvalSymlinks(cwd)
	if err == nil {
		if profile, exists := workspace.Projects[absPath]; exists {
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
