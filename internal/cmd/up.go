package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with creating and starting containers
func UpCmd() *cobra.Command {
	// Up container options
	var opts containers.UpOptions

	// This command and its subcommands (if i add any) would orchestrate containers handling
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Create and start containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return containers.UpContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by UpOptions
	// Flag for passing filepath
	upCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", "", "Path to the configuration file")

	return upCmd
}
