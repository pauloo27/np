package main

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"code.db.cafe/pauloo27/np/config"
	"github.com/spf13/cobra"
)

func newMetaCmd() *cobra.Command {
	metaCmd := &cobra.Command{
		Use:   "meta",
		Short: "Manage metadata for the current directory",
	}

	metaCmd.AddCommand(newMetaSetCmd())
	metaCmd.AddCommand(newMetaUnsetCmd())
	metaCmd.AddCommand(newMetaGetCmd())

	return metaCmd
}

func newMetaSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key=value...]",
		Short: "Add or override metadata keys",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pairs, err := parseMetaPairs(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			project, cwd := mustGetCurrentProject()

			if project.Metadata == nil {
				project.Metadata = make(map[string]string, len(pairs))
			}
			maps.Copy(project.Metadata, pairs)

			mustSaveWorkspace()
			fmt.Printf("Set metadata for %s\n", cwd)
		},
	}
}

func newMetaUnsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset [key...]",
		Short: "Remove metadata keys",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			project, cwd := mustGetCurrentProject()

			for _, key := range args {
				if _, ok := project.Metadata[key]; !ok {
					fmt.Fprintf(os.Stderr, "metadata key '%s' not set for %s\n", key, cwd)
					os.Exit(1)
				}
			}

			for _, key := range args {
				delete(project.Metadata, key)
			}

			mustSaveWorkspace()
			fmt.Printf("Unset metadata for %s\n", cwd)
		},
	}
}

func newMetaGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Show one metadata value, or all keys",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			project, cwd := mustGetCurrentProject()

			if len(args) == 1 {
				value, ok := project.Metadata[args[0]]
				if !ok {
					fmt.Fprintf(os.Stderr, "metadata key '%s' not set for %s\n", args[0], cwd)
					os.Exit(1)
				}
				fmt.Println(value)
				return
			}

			for _, key := range slices.Sorted(maps.Keys(project.Metadata)) {
				fmt.Printf("%s=%s\n", key, project.Metadata[key])
			}
		},
	}
}

func mustGetCurrentProject() (*config.Project, string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current directory: %v\n", err)
		os.Exit(1)
	}

	project, ok := workspace.Projects[cwd]
	if !ok {
		fmt.Fprintf(os.Stderr, "no profile set for %s\n", cwd)
		os.Exit(1)
	}

	return project, cwd
}

func mustSaveWorkspace() {
	if err := workspace.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "error saving workspace: %v\n", err)
		os.Exit(1)
	}
}
