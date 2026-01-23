package main

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "np",
	Short: "Nix project development environment manager",
	Run: func(cmd *cobra.Command, args []string) {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		nixArgs := []string{"nix", "develop", "-c", shell}
		env := append(os.Environ(), "USING_NIX_DEV=local")

		if err := syscall.Exec("/usr/bin/nix", nixArgs, env); err != nil {
			panic(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
