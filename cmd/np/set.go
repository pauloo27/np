package main

import (
	"fmt"
	"os"
	"slices"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

func newSetCmd() *cobra.Command {
	return &cobra.Command{
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

			project := config.NewProject(profile)

			workspace.Projects[cwd] = project

			if err := workspace.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Set profile '%s' for %s\n", profile, cwd)
		},
	}
}
