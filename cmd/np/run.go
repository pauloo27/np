package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [profile]",
	Short: "Run nix develop shell",
	Args:  cobra.MaximumNArgs(1),
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
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error getting home directory: %v\n", err)
			os.Exit(1)
		}

		nixDevProfilesPath := filepath.Join(homeDir, nixDevProfilesDir)

		var profile string
		var useLocalFlake bool

		if len(args) > 0 {
			profile = args[0]
		} else {
			detectedProfile, found := determineProfileName()
			if found {
				profile = detectedProfile
				if profile == "local" {
					useLocalFlake = true
				}
			} else {
				fmt.Fprintf(os.Stderr, "no profile specified and could not determine one automatically\n")
				fmt.Fprintf(os.Stderr, "use 'np set <profile>' to set a profile for this directory\n")
				listAvailableProfiles(nixDevProfilesPath)
				os.Exit(1)
			}
		}

		var nixArgs []string
		var env []string

		if useLocalFlake {
			nixArgs = []string{"nix", "develop", "-c", shell}
			env = append(os.Environ(), "USING_NIX_DEV=local")
		} else {
			profilePath := filepath.Join(nixDevProfilesPath, profile)

			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "profile '%s' not found\n", profile)
				listAvailableProfiles(nixDevProfilesPath)
				os.Exit(1)
			}

			nixArgs = []string{"nix", "develop", profilePath, "-c", shell}
			env = append(os.Environ(), fmt.Sprintf("USING_NIX_DEV=%s", profile))
		}

		if err := syscall.Exec("/usr/bin/nix", nixArgs, env); err != nil {
			panic(err)
		}
	},
}
