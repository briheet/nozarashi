package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StatusCmd(ctx context.Context) *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Checking in the status if the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := statusSystemContainers(ctx)
			return err
		},
	}
	return statusCmd
}

func statusSystemContainers(ctx context.Context) error {
	// Build base command
	systemStatusCmd := exec.CommandContext(
		ctx,
		containers.ContainersCliName,
		containers.ContainersSystemStatusArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStatusCmd.Stderr = os.Stderr
	systemStatusCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStatusCmd.Run()
}
