package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
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

		nixArgs = append(nixArgs, "-c", shell)

		if err := syscall.Exec("/usr/bin/nix", nixArgs, os.Environ()); err != nil {
			panic(err)
		}
	},
}
