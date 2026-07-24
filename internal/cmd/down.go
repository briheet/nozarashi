package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with stopping containers
func DownCmd() *cobra.Command {
	// Down container options
	opts := containers.ContainerOptions{
		FilePath: "nozarashi.toml",
	}

	// This command stops containers from the configured project.
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Stop running containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return containers.DownContainers(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path.
	downCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file")

	return downCmd
}
