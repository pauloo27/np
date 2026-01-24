package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	windowCountFlag int
)

var tmuxCmd = &cobra.Command{
	Use:               "tmux [profile]",
	Short:             "Start a tmux session with nix develop shell",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: profileCompletion(true),
	Run: func(cmd *cobra.Command, args []string) {
		profilesPath := getProfilesPath()

		profile, useLocalFlake, ok := resolveProfile(args, profilesPath)
		if !ok {
			os.Exit(1)
		}

		windowCount := windowCountFlag
		if windowCount <= 0 {
			windowCount = config.Tmux.WindowCount
			if windowCount <= 0 {
				windowCount = 1
			}
		}

		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		var nixCmd string
		if useLocalFlake {
			nixCmd = fmt.Sprintf("nix develop -c %s", shell)
		} else {
			profilePath := filepath.Join(profilesPath, profile)
			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "profile '%s' not found\n", profile)
				listAvailableProfiles(profilesPath)
				os.Exit(1)
			}
			nixCmd = fmt.Sprintf("nix develop %s -c %s", profilePath, shell)
		}

		cwd, _ := os.Getwd()
		sessionName := filepath.Base(cwd)

		checkSession := exec.Command("tmux", "has-session", "-t", sessionName)
		if checkSession.Run() == nil {
			tmuxAttach := []string{"tmux", "attach-session", "-t", sessionName}
			if err := runCommand(tmuxAttach...); err != nil {
				fmt.Fprintf(os.Stderr, "error attaching to tmux session: %v\n", err)
				os.Exit(1)
			}
			return
		}

		tmuxNewSession := []string{"tmux", "new-session", "-d", "-s", sessionName, "-c", cwd, nixCmd}
		if err := runCommand(tmuxNewSession...); err != nil {
			fmt.Fprintf(os.Stderr, "error creating tmux session: %v\n", err)
			os.Exit(1)
		}

		for i := 1; i < windowCount; i++ {
			tmuxNewWindow := []string{"tmux", "new-window", "-t", sessionName, "-c", cwd, nixCmd}
			if err := runCommand(tmuxNewWindow...); err != nil {
				fmt.Fprintf(os.Stderr, "error creating window %d: %v\n", i+1, err)
			}
		}

		tmuxAttach := []string{"tmux", "attach-session", "-t", sessionName}
		if err := runCommand(tmuxAttach...); err != nil {
			fmt.Fprintf(os.Stderr, "error attaching to tmux session: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	tmuxCmd.Flags().IntVarP(&windowCountFlag, "count", "c", 0, "Number of tmux windows to create")
}
