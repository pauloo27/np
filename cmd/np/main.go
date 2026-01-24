package main

import (
	"os"

	"github.com/spf13/cobra"
)

const nixDevProfilesDir = ".config/nix-conf/dev"

var rootCmd = &cobra.Command{
	Use:   "np",
	Short: "Nix project development environment manager",
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(tmuxCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
