package system

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StopCmd() *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Args:  cobra.NoArgs,
		Short: "Stopping all the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StopSystemContainers(cmd.Context())
			return err
		},
	}
	return stopCmd
}
