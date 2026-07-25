package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with listing running containers
func PsCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	// This command and its subcommands (if i add any) would orchestrate containers handling
	psCmd := &cobra.Command{
		Use:   "ps",
		Short: "List running containers.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return containers.PsContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by ContainerOptions
	// Flag for passing filepath
	psCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	// Flag for listing every container in the runtime
	psCmd.Flags().BoolVarP(&opts.All, "all", "a", false, "List all containers")

	return psCmd
}
