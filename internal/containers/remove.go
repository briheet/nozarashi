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

// Removes project containers in reverse dependency order.
func removeContainers(ctx context.Context, graph *specs.Graph) error {
	for _, serviceNode := range slices.Backward(graph.Nodes) {
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := replicas; replica >= 1; replica-- {
			containerName := fmt.Sprintf("%s-%d", serviceNode.ServiceName, replica)

			// Skip service containers that have not been created.
			inspectContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerInspectArgs(containerName)...,
			)
			if err := inspectContainerCmd.Run(); err != nil {
				continue
			}

			deleteContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerDeleteArgs(containerName)...,
			)
			deleteContainerCmd.Stderr = os.Stderr
			deleteContainerCmd.Stdout = os.Stdout

			if err := deleteContainerCmd.Run(); err != nil {
				return fmt.Errorf("delete service container %q: %w", containerName, err)
			}
		}
	}

	return nil
}

// Removes project-scoped images before rebuilding them.
func removeImages(ctx context.Context, graph *specs.Graph) error {
	for _, serviceNode := range slices.Backward(graph.Nodes) {
		// Registry images may be shared by containers from other projects.
		if serviceNode.Spec.Type == specs.ServiceTypeOCI {
			continue
		}

		// Skip project images that have not been built.
		inspectImageCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerImageInspectArgs(serviceNode.ServiceName)...,
		)
		if err := inspectImageCmd.Run(); err != nil {
			continue
		}

		deleteImageCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerImageDeleteArgs(serviceNode.ServiceName)...,
		)
		deleteImageCmd.Stderr = os.Stderr
		deleteImageCmd.Stdout = os.Stdout

		if err := deleteImageCmd.Run(); err != nil {
			return fmt.Errorf("delete service image %q: %w", serviceNode.ServiceName, err)
		}
	}

	return nil
}

// Removes project networks.
func removeNetworks(ctx context.Context, graph *specs.Graph) error {
	for _, resource := range graph.Resources {
		if resource.Kind != specs.ResourceNetwork {
			continue
		}

		// Skip project networks that have not been created.
		inspectNetworkCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerNetworkInspectArgs(resource.Name)...,
		)
		if err := inspectNetworkCmd.Run(); err != nil {
			continue
		}

		deleteNetworkCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerNetworkDeleteArgs(resource.Name)...,
		)
		deleteNetworkCmd.Stderr = os.Stderr
		deleteNetworkCmd.Stdout = os.Stdout

		if err := deleteNetworkCmd.Run(); err != nil {
			return fmt.Errorf("delete network %q: %w", resource.Name, err)
		}
	}

	return nil
}

// Removes project volumes.
func removeVolumes(ctx context.Context, graph *specs.Graph) error {
	for _, resource := range graph.Resources {
		if resource.Kind != specs.ResourceVolume {
			continue
		}

		// Skip project volumes that have not been created.
		inspectVolumeCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerVolumeInspectArgs(resource.Name)...,
		)
		if err := inspectVolumeCmd.Run(); err != nil {
			continue
		}

		deleteVolumeCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerVolumeDeleteArgs(resource.Name)...,
		)
		deleteVolumeCmd.Stderr = os.Stderr
		deleteVolumeCmd.Stdout = os.Stdout

		if err := deleteVolumeCmd.Run(); err != nil {
			return fmt.Errorf("delete volume %q: %w", resource.Name, err)
		}
	}

	return nil
}

// Removes every service image declared by the project.
func destroyImages(ctx context.Context, graph *specs.Graph) error {
	for _, serviceNode := range slices.Backward(graph.Nodes) {
		imageName := serviceNode.ServiceName
		if serviceNode.Spec.Type == specs.ServiceTypeOCI {
			imageName = serviceNode.Spec.Reference
		}

		// Skip service images that are not present.
		inspectImageCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerImageInspectArgs(imageName)...,
		)
		if err := inspectImageCmd.Run(); err != nil {
			continue
		}

		deleteImageCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerImageDeleteArgs(imageName)...,
		)
		deleteImageCmd.Stderr = os.Stderr
		deleteImageCmd.Stdout = os.Stdout

		if err := deleteImageCmd.Run(); err != nil {
			return fmt.Errorf("delete service image %q: %w", imageName, err)
		}
	}

	return nil
}
