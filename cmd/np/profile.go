package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

func newProfileCmd() *cobra.Command {
	var variation string

	cmd := &cobra.Command{
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

			if variation == "" {
				variation = resolveVariation()
			}

			var nixArgs []string
			var env []string

			if useLocalFlake {
				flakeRef := "."
				if variation != "" {
					flakeRef = ".#" + variation
				}
				nixArgs = []string{"nix", "develop", flakeRef, "-c", shell}
				env = append(os.Environ(), "USING_NIX_DEV=local")
			} else {
				profilePath := filepath.Join(profilesPath, profile)

				if _, err := os.Stat(profilePath); os.IsNotExist(err) {
					fmt.Fprintf(os.Stderr, "profile '%s' not found\n", profile)
					listAvailableProfiles(profilesPath)
					os.Exit(1)
				}

				flakeRef := profilePath
				if variation != "" {
					flakeRef = profilePath + "#" + variation
				}
				nixArgs = []string{"nix", "develop", flakeRef, "-c", shell}
				env = append(os.Environ(), fmt.Sprintf("USING_NIX_DEV=%s", profile))
			}

			nixPath, err := getBinPath("nix")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := syscall.Exec(nixPath, nixArgs, env); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&variation, "variation", "v", "", "Nix develop variation (e.g. node22)")

	return cmd
}
