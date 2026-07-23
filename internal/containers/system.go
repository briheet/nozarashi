package containers

import (
	"context"
	"os"
	"os/exec"
)

// This wraps over apple's container cli for starting it
func StartSystemContainers(ctx context.Context) error {
	// Build base command
	systemStartCmd := exec.CommandContext(
		ctx,
		ContainersCliName,
		ContainersSystemStartArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStartCmd.Stderr = os.Stderr
	systemStartCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStartCmd.Run()
}

// This wraps over apple's container cli for getting its status
func StatusSystemContainers(ctx context.Context) error {
	// Build base command
	systemStatusCmd := exec.CommandContext(
		ctx,
		ContainersCliName,
		ContainersSystemStatusArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStatusCmd.Stderr = os.Stderr
	systemStatusCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStatusCmd.Run()
}

// This wraps over apple's container cli for stopping it
func StopSystemContainers(ctx context.Context) error {
	// Build base command
	systemStopCmd := exec.CommandContext(
		ctx,
		ContainersCliName,
		ContainersSystemStopArgs...,
	)

	// Point the Command Error and Output to Standard Error and Output
	systemStopCmd.Stderr = os.Stderr
	systemStopCmd.Stdout = os.Stdout

	// Run and return if any issues
	return systemStopCmd.Run()
}
