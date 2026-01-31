package main

import (
	"os"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

var (
	// global state, fuck it, we ball
	cfg       *config.Config
	workspace *config.Workspace
)

var rootCmd = &cobra.Command{
	Use:   "np",
	Short: "Nix project development environment manager",
}

func init() {
	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newProfileCmd())
	rootCmd.AddCommand(newSetCmd())
	rootCmd.AddCommand(newDevCmd())
	rootCmd.AddCommand(newTmuxCmd())
	rootCmd.AddCommand(newShellCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newLoadCmd())
}

func main() {
	cfg, _ = config.LoadConfig()
	if cfg == nil {
		cfg = &config.Config{}
	}

	workspace, _ = config.LoadWorkspace(cfg)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
