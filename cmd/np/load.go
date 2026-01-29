package main

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load the profile if it should be using one and isn't",
	Run: func(cmd *cobra.Command, args []string) {
		shouldUse := os.Getenv("SHOULD_USE_NIX_DEV")
		currentlyUsing := os.Getenv("USING_NIX_DEV")

		if shouldUse == "" || currentlyUsing != "" {
			return
		}

		if currentlyUsing == shouldUse {
			return
		}

		execArgs := []string{os.Args[0], "profile", shouldUse}
		if err := syscall.Exec(os.Args[0], execArgs, os.Environ()); err != nil {
			panic(err)
		}
	},
}
