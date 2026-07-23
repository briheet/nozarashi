package system

import (
	"context"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StatusCmd(ctx context.Context) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Checking in the status if the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StatusSystemContainers(ctx)
			return err
		},
	}
	return statusCmd
}
