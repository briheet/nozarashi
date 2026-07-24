package containers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/briheet/nozarashi/internal/specs"
)

// Stops service containers in reverse dependency order.
func stopContainers(ctx context.Context, graph *specs.Graph) error {
	for _, serviceNode := range slices.Backward(graph.Nodes) {
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := replicas; replica >= 1; replica-- {
			containerName := fmt.Sprintf("%s-%d", serviceNode.ServiceName, replica)

			// Inspect containers running and report back.
			// Skip service containers that have not been created.
			inspectContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerInspectArgs(containerName)...,
			)
			if err := inspectContainerCmd.Run(); err != nil {
				continue
			}

			// Stop both running and already stopped containers safely.
			stopContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerStopArgs(containerName)...,
			)
			stopContainerCmd.Stderr = os.Stderr
			stopContainerCmd.Stdout = os.Stdout

			if err := stopContainerCmd.Run(); err != nil {
				return fmt.Errorf("stop service container %q: %w", containerName, err)
			}
		}
	}

	return nil
}
