package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command displays details for project service containers.
func InspectCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	inspectCmd := &cobra.Command{
		Use:   "inspect <services...>",
		Short: "Display service container details",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return containers.InspectContainers(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path
	inspectCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return inspectCmd
}
