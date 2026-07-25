package cmd

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

// This command executes a command in a running service container.
func ExecCmd() *cobra.Command {
	// Container options
	opts := containers.ContainerOptions{
		// Default config filepath
		FilePath: "nozarashi.toml",

		// Default service replica
		Replica: 1,
	}

	execCmd := &cobra.Command{
		Use:   "exec <service> <command> [args...]",
		Short: "Execute a command in a service container",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args
			return containers.ExecContainers(cmd.Context(), opts)
		},
	}

	// Flag for passing the configuration file path
	execCmd.Flags().StringVarP(&opts.FilePath, "filepath", "f", opts.FilePath, "Path to the configuration file.")

	// Flag for selecting a service replica
	execCmd.Flags().IntVarP(&opts.Replica, "replica", "r", opts.Replica, "Service replica to execute in")

	// Flags for interactive terminal execution
	execCmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "Keep standard input open")
	execCmd.Flags().BoolVarP(&opts.TTY, "tty", "t", false, "Allocate a terminal")

	// Treat every flag after the service name as part of the executed command.
	execCmd.Flags().SetInterspersed(false)

	return execCmd
}
