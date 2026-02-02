package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newUnsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset",
		Short: "Unset workspace for the current directory",
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error getting current directory: %v\n", err)
				os.Exit(1)
			}

			if _, ok := workspace.Projects[cwd]; !ok {
				fmt.Fprintf(os.Stderr, "no profile set for %s\n", cwd)
				os.Exit(1)
			}

			delete(workspace.Projects, cwd)

			if err := workspace.Save(); err != nil {
				fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Unset profile for %s\n", cwd)
		},
	}
}
