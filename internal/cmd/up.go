package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

// This command usually deals with up and starting containers
func UpCmd(ctx context.Context) *cobra.Command {
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Create and start containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return upCmd
}
