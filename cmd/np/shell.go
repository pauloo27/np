package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

func newShellCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shell [packages...]",
		Short: "Run nix shell with packages from nixpkgs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/sh"
			}

			nixArgs := []string{"nix", "shell"}

			for _, pkg := range args {
				nixArgs = append(nixArgs, fmt.Sprintf("nixpkgs#%s", pkg))
			}

			nixPath, err := getBinPath("nix")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			nixArgs = append(nixArgs, "-c", shell)

			if err := syscall.Exec(nixPath, nixArgs, os.Environ()); err != nil {
				panic(err)
			}
		},
	}
}
