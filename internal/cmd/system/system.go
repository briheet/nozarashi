package system

import (
	"github.com/spf13/cobra"
)

func SystemCmd() *cobra.Command {
	systemCmd := &cobra.Command{
		Use:   "system",
		Args:  cobra.NoArgs,
		Short: "Manage the container system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	systemCmd.AddCommand(StartCmd())
	systemCmd.AddCommand(StatusCmd())
	systemCmd.AddCommand(StopCmd())

	return systemCmd
}
