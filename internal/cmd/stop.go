package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/containers"
	"github.com/spf13/cobra"
)

func StopCmd(ctx context.Context) *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stopping all the containers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := stopSystemContainers(ctx)
			return err
		},
	}
	return stopCmd
}

func stopSystemContainers(ctx context.Context) error {
	// Build base command
	systemStopCmd := exec.CommandContext(
		ctx,
		containers.ContainersCliName,
		containers.ContainersSystemStopArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStopCmd.Stderr = os.Stderr
	systemStopCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStopCmd.Run()
}
