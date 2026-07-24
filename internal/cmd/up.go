package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with creating and starting containers
func UpCmd() *cobra.Command {
	// Up container options
	opts := containers.ContainerOptions{
		FilePath: "nozarashi.toml",
	}

	// This command and its subcommands (if i add any) would orchestrate containers handling
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Create and start containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Add args
			opts.Args = args
			return containers.UpContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by UpOptions
	// Flag for passing filepath
	upCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	// Flag for rebuilding images, else skip
	upCmd.Flags().BoolVarP(&opts.Build, "build", "b", opts.Build, "Pass for rebuilding of images.")

	return upCmd
}
