package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// global state, fuck it, we ball
	config    *Config
	workspace *Workspace
)

var rootCmd = &cobra.Command{
	Use:   "np",
	Short: "Nix project development environment manager",
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(tmuxCmd)
	rootCmd.AddCommand(xCmd)
}

func main() {
	config, _ = loadConfig()
	if config == nil {
		config = &Config{}
	}

	workspace, _ = loadWorkspace()
	if workspace == nil {
		workspace = &Workspace{Projects: make(map[string]string)}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
