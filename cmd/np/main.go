package main

import (
	"fmt"
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
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid config: %v\n", err)
		os.Exit(1)
	}

	workspace, err = config.LoadWorkspace(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid workspace: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
