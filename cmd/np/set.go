package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [profile]",
	Short: "Set the profile for the current directory",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		profilesPath := filepath.Join(homeDir, nixDevProfilesDir)
		profiles, err := getAvailableProfiles(profilesPath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return profiles, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting current directory: %v\n", err)
			os.Exit(1)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(homeDir, nixDevProfilesDir, "config.toml")

		config, err := loadConfig()
		if err != nil {
			config = &Config{Projects: make(map[string]string)}
		}

		config.Projects[cwd] = profile

		configDir := filepath.Dir(configPath)
		if err := os.MkdirAll(configDir, 0750); err != nil {
			fmt.Fprintf(os.Stderr, "error creating config directory: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		if err := toml.NewEncoder(file).Encode(config); err != nil {
			fmt.Fprintf(os.Stderr, "error writing config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set profile '%s' for %s\n", profile, cwd)
	},
}
