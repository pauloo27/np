package main

import (
	"fmt"
	"os"
	"path/filepath"

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

		workspace, err := loadWorkspace()
		if err != nil {
			workspace = &Workspace{Projects: make(map[string]string)}
		}

		workspace.Projects[cwd] = profile

		if err := saveWorkspace(workspace); err != nil {
			fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set profile '%s' for %s\n", profile, cwd)
	},
}
