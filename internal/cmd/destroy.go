package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command destroys all runtime objects declared by a project.
func DestroyCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Delete project containers, networks, volumes and images",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return containers.DestroyContainers(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path
	destroyCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return destroyCmd
}
