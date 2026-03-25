package main

import (
	"fmt"
	"os"
	"slices"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

func newSetCmd() *cobra.Command {
	var (
		windowCountFlag    int
		windowsCommandFlag []string
		sessionName        string
	)

	setCmd := &cobra.Command{
		Use:               "set [profile]",
		Short:             "Set the profile for the current directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: profileCompletion(true),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]

			if windowCountFlag < 0 {
				fmt.Fprintf(os.Stderr, "tmux window count cannot be negative \n")
				os.Exit(1)
			}

			profilesPath := getProfilesPath()
			availableProfiles, err := getAvailableProfiles(profilesPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading profiles: %v\n", err)
				os.Exit(1)
			}

			profileExists := slices.Contains(availableProfiles, profile)

			if !profileExists && profile != "local" && profile != "none" {
				fmt.Fprintf(os.Stderr, "profile '%s' does not exist\n", profile)
				listAvailableProfiles(profilesPath)
				os.Exit(1)
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting current directory: %v\n", err)
				os.Exit(1)
			}

			windows := buildTmuxWindows(windowCountFlag, windowsCommandFlag)

			project := config.NewProject(profile, windows, sessionName)

			workspace.Projects[cwd] = project

			if err := workspace.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Set profile '%s' for %s\n", profile, cwd)
		},
	}

	setCmd.Flags().IntVarP(&windowCountFlag, "count", "c", 0, "Number of tmux windows for the project")
	setCmd.Flags().StringArrayVarP(&windowsCommandFlag, "window", "w", []string{}, "Add a window with a command")
	setCmd.Flags().StringVarP(&sessionName, "name", "n", "", "Set tmux session name")

	return setCmd
}
