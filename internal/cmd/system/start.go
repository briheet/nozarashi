package system

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Args:  cobra.NoArgs,
		Short: "Starting the apple containers to work with",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StartSystemContainers(cmd.Context())
			return err
		},
	}
	return startCmd
}
