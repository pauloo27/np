package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Run nix develop shell in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		shell := getShell()

		nixArgs := []string{"nix", "develop", "-c", shell}

		nixPath, err := getBinPath("nix")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		if err := syscall.Exec(nixPath, nixArgs, os.Environ()); err != nil {
			panic(err)
		}
	},
}
