package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	buildCmd := &cobra.Command{
		Use:   "build [services...]",
		Short: "Build images, remove existing containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return containers.BuildContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by ContainerOptions
	// Flag for passing filepath
	buildCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	return buildCmd
}
