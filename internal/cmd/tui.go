package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/briheet/nozarashi/internal/tui/view"
	"github.com/spf13/cobra"
)

// Spins up a tui for project specific containers running
func TuiCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	// This command helps to spin up a tui for workflow.
	tuiCmd := &cobra.Command{
		Use:   "tui",
		Args:  cobra.NoArgs,
		Short: "Tui viewer for seeing container logs and other specifics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return view.Run(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by ContainerOptions
	// Flag for passing filepath
	tuiCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return tuiCmd
}
