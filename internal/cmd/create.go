package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command creates volumes and networks declared by the project.
func CreateCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create project volumes and networks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return containers.CreateResources(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path
	createCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return createCmd
}
