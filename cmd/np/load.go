package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

func newLoadCmd() *cobra.Command {
	return &cobra.Command{
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

			var npPath string
			var err error

			// that should match "./np", "/bin/np" etc
			if strings.Contains(npPath, "/") {
				npPath = os.Args[0]
			} else {
				npPath, err = getBinPath("np")
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
					os.Exit(1)
				}
			}

			execArgs := []string{npPath, "profile", shouldUse}
			if err := syscall.Exec(npPath, execArgs, os.Environ()); err != nil {
				panic(err)
			}
		},
	}
}
