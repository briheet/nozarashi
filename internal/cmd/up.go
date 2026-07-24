package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with creating and starting containers
func UpCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	// This command and its subcommands (if i add any) would orchestrate containers handling
	upCmd := &cobra.Command{
		Use:   "up [services...]",
		Short: "Create and start containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Add args
			opts.Args = args
			return containers.UpContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by ContainerOptions
	// Flag for passing filepath
	upCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return upCmd
}
