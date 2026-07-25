package containers

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/specs"
)

// Helps in creating Resources
func createResources(ctx context.Context, graph *specs.Graph) error {
	// Create volumes defined in the graph
	if err := createVolumes(ctx, graph); err != nil {
		return err
	}

	// Create networks defined in the graph
	if err := createNetworks(ctx, graph); err != nil {
		return err
	}

	return nil
}

// Creates and starts service containers in dependency order.
func createContainers(ctx context.Context, graph *specs.Graph) error {
	for _, serviceNode := range graph.Nodes {
		// OCI services use their registry reference. Built services use their local image tag.
		imageName := serviceNode.ServiceName
		if serviceNode.Spec.Type == specs.ServiceTypeOCI {
			imageName = serviceNode.Spec.Reference
		}

		// An omitted replica count creates one container.
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := 1; replica <= replicas; replica++ {
			containerName := fmt.Sprintf("%s-%d", serviceNode.ServiceName, replica)

			// Start an existing container without recreating its resources.
			inspectContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerInspectArgs(containerName)...,
			)
			if err := inspectContainerCmd.Run(); err == nil {
				startContainerCmd := exec.CommandContext(
					ctx,
					ContainerCliName,
					ContainerStartArgs(containerName)...,
				)
				startContainerCmd.Stderr = os.Stderr
				startContainerCmd.Stdout = os.Stdout

				if err := startContainerCmd.Run(); err != nil {
					return fmt.Errorf("start service container %q: %w", containerName, err)
				}
				continue
			}

			// Create and start a new detached container.
			runContainerCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerRunArgs(
					graph.Project.Project.Name,
					containerName,
					imageName,
					fmt.Sprintf(
						"%s/%s",
						specs.CurrentPlatform.OperatingSystem,
						specs.CurrentPlatform.Architecture,
					),
					serviceNode.Spec,
				)...,
			)
			runContainerCmd.Stderr = os.Stderr
			runContainerCmd.Stdout = os.Stdout

			if err := runContainerCmd.Run(); err != nil {
				return fmt.Errorf("create service container %q: %w", containerName, err)
			}
		}
	}

	return nil
}

// Wraps up container specifics for volume
func createVolumes(ctx context.Context, graph *specs.Graph) error {
	for _, resource := range graph.Resources {
		if resource.Kind != specs.ResourceVolume {
			continue
		}

		// Reuse an existing project volume.
		inspectVolumeCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerVolumeInspectArgs(resource.Name)...,
		)
		if err := inspectVolumeCmd.Run(); err == nil {
			continue
		}

		// Create the volume with Apple Container's native volume backend.
		createVolumeCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerVolumeCreateArgs(resource.Name)...,
		)
		createVolumeCmd.Stderr = os.Stderr
		createVolumeCmd.Stdout = os.Stdout

		if err := createVolumeCmd.Run(); err != nil {
			return fmt.Errorf("create volume %q: %w", resource.Name, err)
		}
	}

	return nil
}

// Wraps up container specifics for networks
func createNetworks(ctx context.Context, graph *specs.Graph) error {
	for _, resource := range graph.Resources {
		if resource.Kind != specs.ResourceNetwork {
			continue
		}

		// Reuse an existing project network.
		inspectNetworkCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerNetworkInspectArgs(resource.Name)...,
		)
		if err := inspectNetworkCmd.Run(); err == nil {
			continue
		}

		// Create the network with Apple Container's native network backend.
		createNetworkCmd := exec.CommandContext(
			ctx,
			ContainerCliName,
			ContainerNetworkCreateArgs(resource.Name)...,
		)
		createNetworkCmd.Stderr = os.Stderr
		createNetworkCmd.Stdout = os.Stdout

		if err := createNetworkCmd.Run(); err != nil {
			return fmt.Errorf("create network %q: %w", resource.Name, err)
		}
	}

	return nil
}
