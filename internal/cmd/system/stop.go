package system

import (
	"context"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StopCmd(ctx context.Context) *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stopping all the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StopSystemContainers(ctx)
			return err
		},
	}
	return stopCmd
}
