package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command usually deals with getting logs of containers.
func LogsCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",
	}

	logsCmd := &cobra.Command{
		Use:   "logs [services...]",
		Short: "Get logs for containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return containers.LogsContainers(cmd.Context(), opts)
		},
	}

	// Add required flags from here. This will be picked by ContainerOptions
	// Flag for passing filepath
	logsCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	// Flag for passing the number of lines of logs we want
	logsCmd.Flags().IntVarP(&opts.Number, "number", "n", 0, "Number of lines of logs to print")

	return logsCmd
}
