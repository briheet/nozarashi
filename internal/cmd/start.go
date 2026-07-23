package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StartCmd(ctx context.Context) *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Starting the apple containers to work with",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := startSystemContainers(ctx)
			return err
		},
	}
	return startCmd
}

func startSystemContainers(ctx context.Context) error {
	// Build base command
	systemStartCmd := exec.CommandContext(
		ctx,
		containers.ContainersCliName,
		containers.ContainersSystemStartArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStartCmd.Stderr = os.Stderr
	systemStartCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStartCmd.Run()
}
