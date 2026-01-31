package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			profilesPath := getProfilesPath()
			profiles, err := getAvailableProfiles(profilesPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading profiles: %v\n", err)
				os.Exit(1)
			}

			if len(profiles) == 0 {
				fmt.Println("no profiles found")
				return
			}

			fmt.Println("available profiles:")
			for _, profile := range profiles {
				fmt.Printf("  %s\n", profile)
			}
		},
	}
}
