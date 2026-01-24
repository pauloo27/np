package main

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run nix develop shell in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		nixArgs := []string{"nix", "develop", "-c", shell}

		if err := syscall.Exec("/usr/bin/nix", nixArgs, os.Environ()); err != nil {
			panic(err)
		}
	},
}
