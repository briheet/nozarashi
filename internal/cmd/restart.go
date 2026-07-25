package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command restarts project service containers.
func RestartCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	restartCmd := &cobra.Command{
		Use:   "restart [services...]",
		Short: "Restart service containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return containers.RestartContainers(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path
	restartCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return restartCmd
}
