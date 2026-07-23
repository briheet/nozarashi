package system

import (
	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Checking in the status if the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := containers.StatusSystemContainers(cmd.Context())
			return err
		},
	}
	return statusCmd
}
