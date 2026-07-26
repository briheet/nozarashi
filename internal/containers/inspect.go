package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/render"
	"github.com/briheet/nozarashi/internal/specs"
)

// This function inspects containers for selected project services.
func InspectContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running.
	if err := StatusSystemContainers(ctx, io.Discard); err != nil {
		return err
	}

	if len(opts.Args) == 0 {
		return fmt.Errorf("inspect requires at least one service")
	}

	projectSpecs, err := config.ParseTOMLConfig(ctx, opts.FilePath)
	if err != nil {
		return err
	}

	if err := validateDependencyGraph(projectSpecs); err != nil {
		return err
	}

	graph, err := buildDependencyGraph(projectSpecs)
	if err != nil {
		return err
	}

	if err := selectServiceNodes(graph, opts.Args); err != nil {
		return err
	}

	containers, err := inspectContainers(ctx, graph)
	if err != nil {
		return err
	}

	if err := render.RenderContainerDetails(ctx, containers); err != nil {
		return fmt.Errorf("render container details: %w", err)
	}

	return nil
}

// Gets detailed configuration and status for containers.
func inspectContainers(ctx context.Context, graph *specs.Graph) (*specs.Containers, error) {
	// Build container names
	containerNames := make([]string, 0, len(graph.Nodes))
	for _, serviceNode := range graph.Nodes {
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := 1; replica <= replicas; replica++ {
			containerNames = append(
				containerNames,
				serviceContainerName(
					graph.Project.Project.Name,
					serviceNode.ServiceName,
					replica,
				),
			)
		}
	}

	// Build inspect command
	inspectCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerInspectArgs(containerNames...)...,
	)
	inspectCmd.Stderr = os.Stderr

	// Catching the output so that can be sent to render
	output, err := inspectCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("inspect containers: %w", err)
	}

	var inspectedContainers []specs.Container
	if err := json.Unmarshal(output, &inspectedContainers); err != nil {
		return nil, fmt.Errorf("decode inspected containers: %w", err)
	}

	return &specs.Containers{Items: inspectedContainers}, nil
}
