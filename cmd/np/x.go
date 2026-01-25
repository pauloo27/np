package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

var xCmd = &cobra.Command{
	Use:   "x [package] [-- args...]",
	Short: "Run a package from nixpkgs",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkg := args[0]

		nixArgs := []string{"nix", "run", fmt.Sprintf("nixpkgs#%s", pkg)}

		if len(args) > 1 {
			nixArgs = append(nixArgs, "--")
			nixArgs = append(nixArgs, args[1:]...)
		}

		if err := syscall.Exec("/usr/bin/nix", nixArgs, os.Environ()); err != nil {
			panic(err)
		}
	},
}
