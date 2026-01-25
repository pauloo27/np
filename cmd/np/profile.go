package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:               "profile [profile]",
	Short:             "Run nix develop shell with profile",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: profileCompletion(true),
	Run: func(cmd *cobra.Command, args []string) {
		shell := getShell()

		profilesPath := getProfilesPath()

		profile, useLocalFlake, ok := resolveProfile(args, profilesPath)
		if !ok {
			os.Exit(1)
		}

		var nixArgs []string
		var env []string

		if useLocalFlake {
			nixArgs = []string{"nix", "develop", "-c", shell}
			env = append(os.Environ(), "USING_NIX_DEV=local")
		} else {
			profilePath := filepath.Join(profilesPath, profile)

			if _, err := os.Stat(profilePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "profile '%s' not found\n", profile)
				listAvailableProfiles(profilesPath)
				os.Exit(1)
			}

			nixArgs = []string{"nix", "develop", profilePath, "-c", shell}
			env = append(os.Environ(), fmt.Sprintf("USING_NIX_DEV=%s", profile))
		}

		nixPath, err := getNixPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		if err := syscall.Exec(nixPath, nixArgs, env); err != nil {
			panic(err)
		}
	},
}
