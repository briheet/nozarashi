package containers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/specs"
)

// This wraps over apple's container cli and helps us create project resources.
func CreateResources(ctx context.Context, opts ContainerOptions) error {
	// Parse the config file
	projectSpecs, err := config.ParseTOMLConfig(ctx, opts.FilePath)
	if err != nil {
		return err
	}

	// Validate dependency graph and find any inconsistencies
	if err := validateDependencyGraph(projectSpecs); err != nil {
		return err
	}

	// Build dependency graph
	graph, err := buildDependencyGraph(projectSpecs)
	if err != nil {
		return err
	}

	// Elevated create commands also register the project DNS domain.
	// Apple containers need sudo privledge for dns
	// https://github.com/apple/container/blob/main/docs/tutorial.md#set-up-a-local-dns-domain-optional
	if err := createProjectDNS(ctx, projectSpecs.Project.Name); err != nil {
		return err
	}

	// Create volumes and networks declared by the project
	if err := createResources(ctx, graph); err != nil {
		return err
	}

	return nil
}

// Creates the project DNS domain when create is running through sudo.
func createProjectDNS(ctx context.Context, domain string) error {
	// Check the invoking user's container system rather than root's.
	statusSystemCmd, err := createContainerCommand(ctx, ContainerSystemStatusArgs...)
	if err != nil {
		return err
	}
	statusSystemCmd.Stderr = os.Stderr
	statusSystemCmd.Stdout = os.Stdout

	if err := statusSystemCmd.Run(); err != nil {
		return fmt.Errorf("check container system status: %w", err)
	}

	if os.Geteuid() != 0 {
		return nil
	}

	listDNSCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerSystemDNSListArgs...,
	)
	output, err := listDNSCmd.Output()
	if err != nil {
		return fmt.Errorf("list container DNS domains: %w", err)
	}

	if slices.Contains(strings.Fields(string(output)), domain) {
		return nil
	}

	createDNSCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerSystemDNSCreateArgs(domain)...,
	)
	createDNSCmd.Stderr = os.Stderr
	createDNSCmd.Stdout = os.Stdout

	if err := createDNSCmd.Run(); err != nil {
		return fmt.Errorf("create container DNS domain %q: %w", domain, err)
	}

	return nil
}

// Runs resource commands as the user who invoked sudo.
func createContainerCommand(ctx context.Context, args ...string) (*exec.Cmd, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if os.Geteuid() != 0 || sudoUser == "" {
		return exec.CommandContext(ctx, ContainerCliName, args...), nil
	}

	containerPath, err := exec.LookPath(ContainerCliName)
	if err != nil {
		return nil, fmt.Errorf("find container cli: %w", err)
	}

	sudoArgs := []string{
		"--user",
		sudoUser,
		"--set-home",
		"--",
		containerPath,
	}
	sudoArgs = append(sudoArgs, args...)

	return exec.CommandContext(ctx, "sudo", sudoArgs...), nil
}

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
		imageName := serviceImageReference(serviceNode)

		// An omitted replica count creates one container.
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := 1; replica <= replicas; replica++ {
			containerName := serviceContainerName(
				graph.Project.Project.Name,
				serviceNode.ServiceName,
				replica,
			)

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
		inspectVolumeCmd, err := createContainerCommand(
			ctx,
			ContainerVolumeInspectArgs(resource.Name)...,
		)
		if err != nil {
			return err
		}
		if err := inspectVolumeCmd.Run(); err == nil {
			continue
		}

		// Create the volume with Apple Container's native volume backend.
		createVolumeCmd, err := createContainerCommand(
			ctx,
			ContainerVolumeCreateArgs(resource.Name)...,
		)
		if err != nil {
			return err
		}
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
		inspectNetworkCmd, err := createContainerCommand(
			ctx,
			ContainerNetworkInspectArgs(resource.Name)...,
		)
		if err != nil {
			return err
		}
		if err := inspectNetworkCmd.Run(); err == nil {
			continue
		}

		// Create the network with Apple Container's native network backend.
		createNetworkCmd, err := createContainerCommand(
			ctx,
			ContainerNetworkCreateArgs(resource.Name)...,
		)
		if err != nil {
			return err
		}
		createNetworkCmd.Stderr = os.Stderr
		createNetworkCmd.Stdout = os.Stdout

		if err := createNetworkCmd.Run(); err != nil {
			return fmt.Errorf("create network %q: %w", resource.Name, err)
		}
	}

	return nil
}
