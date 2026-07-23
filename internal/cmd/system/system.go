package system

import (
	"context"

	"github.com/spf13/cobra"
)

func SystemCmd(ctx context.Context) *cobra.Command {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "Manage the container system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	systemCmd.AddCommand(StartCmd(ctx))
	systemCmd.AddCommand(StatusCmd(ctx))
	systemCmd.AddCommand(StopCmd(ctx))

	return systemCmd
}
