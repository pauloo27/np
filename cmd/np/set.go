package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:               "set [profile]",
	Short:             "Set the profile for the current directory",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: profileCompletion(false),
	Run: func(cmd *cobra.Command, args []string) {
		profile := args[0]

		profilesPath := getProfilesPath()
		availableProfiles, err := getAvailableProfiles(profilesPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading profiles: %v\n", err)
			os.Exit(1)
		}

		profileExists := slices.Contains(availableProfiles, profile)

		if !profileExists {
			fmt.Fprintf(os.Stderr, "profile '%s' does not exist\n", profile)
			listAvailableProfiles(profilesPath)
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting current directory: %v\n", err)
			os.Exit(1)
		}

		workspace.Projects[cwd] = profile

		if err := saveWorkspace(workspace); err != nil {
			fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set profile '%s' for %s\n", profile, cwd)
	},
}
