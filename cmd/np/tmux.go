package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

func newTmuxCmd() *cobra.Command {
	var (
		windowCountFlag    int
		windowsCommandFlag []string
	)

	tmuxCmd := &cobra.Command{
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

			var profileStartCmd string
			if profile != "none" {
				if useLocalFlake {
					profileStartCmd = fmt.Sprintf("%s profile local", os.Args[0])
				} else {
					profileStartCmd = fmt.Sprintf("%s profile %s", os.Args[0], profile)
				}
			}

			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting current work directory: %v", err)
				os.Exit(1)
			}

			project := workspace.Projects[cwd]

			sessionName := filepath.Base(cwd)
			if project != nil && project.Tmux.SessionName != "" {
				sessionName = project.Tmux.SessionName
			}

			checkSession := exec.Command("tmux", "has-session", "-t", "="+sessionName)
			if checkSession.Run() == nil {
				tmuxAttach := []string{"tmux", "attach-session", "-t", "=" + sessionName}
				if err := runCommand(tmuxAttach...); err != nil {
					fmt.Fprintf(os.Stderr, "error attaching to tmux session: %v\n", err)
					os.Exit(1)
				}
				return
			}

			tmuxNewSession := []string{"tmux", "new-session", "-e", "SHOULD_USE_NIX_DEV=" + profile, "-d", "-s", sessionName, "-c", cwd, profileStartCmd}
			if err := runCommand(tmuxNewSession...); err != nil {
				fmt.Fprintf(os.Stderr, "error creating tmux session: %v\n", err)
				os.Exit(1)
			}

			var windows []*config.TmuxWindow

			hasCommandFlags := len(windowsCommandFlag) > 0 || windowCountFlag > 0

			if !hasCommandFlags {
				if project != nil && project.Tmux != nil && len(project.Tmux.Windows) > 0 {
					windows = project.Tmux.Windows
				} else {
					windowCountFlag = 1
				}
			}

			if len(windows) == 0 {
				windows = buildTmuxWindows(windowCountFlag, windowsCommandFlag)
			}

			for i, window := range windows {
				// First window already created by new-session, create the rest
				if i > 0 {
					tmuxNewWindow := []string{"tmux", "new-window", "-t", sessionName, "-c", cwd, profileStartCmd}
					if err := runCommand(tmuxNewWindow...); err != nil {
						fmt.Fprintf(os.Stderr, "error creating additional window: %v\n", err)
					}
				}

				// Send command to window if it has one
				if window.Command != "" {
					windowIndex := cfg.TmuxBaseWindowIndex + i
					target := fmt.Sprintf("%s:%d", sessionName, windowIndex)
					tmuxSendKeys := []string{"tmux", "send-keys", "-t", target, window.Command, "Enter"}
					if err := runCommand(tmuxSendKeys...); err != nil {
						fmt.Fprintf(os.Stderr, "error sending command to window %d: %v\n", windowIndex, err)
					}
				}
			}

			windowSeletion := []string{"tmux", "select-window", "-t", fmt.Sprintf("%d", cfg.TmuxBaseWindowIndex)}
			if err := runCommand(windowSeletion...); err != nil {
				fmt.Fprintf(os.Stderr, "error creating tmux session: %v\n", err)
				os.Exit(1)
			}

			tmuxAttach := []string{"tmux", "attach-session", "-t", sessionName}
			if err := runCommand(tmuxAttach...); err != nil {
				fmt.Fprintf(os.Stderr, "error attaching to tmux session: %v\n", err)
				os.Exit(1)
			}
		},
	}

	tmuxCmd.Flags().IntVarP(&windowCountFlag, "count", "c", 0, "Number of tmux windows to create")
	tmuxCmd.Flags().StringArrayVarP(&windowsCommandFlag, "window", "w", []string{}, "Add a window with a command")

	return tmuxCmd
}
