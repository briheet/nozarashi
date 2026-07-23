package system

import (
	"context"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StartCmd(ctx context.Context) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starting the apple containers to work with",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StartSystemContainers(ctx)
			return err
		},
	}
	return startCmd
}
