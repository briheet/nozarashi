package containers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/render"
	"github.com/briheet/nozarashi/internal/specs"
	"github.com/charmbracelet/x/ansi"
)

// This function wraps over apple's container cli and lists containers.
func PsContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running.
	if err := StatusSystemContainers(ctx, io.Discard); err != nil {
		return err
	}

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

	// Get containers
	containers, err := GetContainers(ctx, opts.All)
	if err != nil {
		return err
	}

	// Select which containers we need to use
	containers = selectProjectContainers(containers, graph, opts.All)

	// Render the data
	if err := render.RenderContainers(ctx, containers); err != nil {
		return fmt.Errorf("render containers: %w", err)
	}

	return nil
}

// GetContainers gets running or all containers from Apple Container.
func GetContainers(ctx context.Context, all bool) (*specs.Containers, error) {
	listCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerListArgs(all)...,
	)
	// Capture stderr because Apple Container writes cursor control codes there.
	stderr := &bytes.Buffer{}
	listCmd.Stderr = stderr

	output, err := listCmd.Output()
	if err != nil {
		// Remove terminal control codes before returning a readable command error.
		message := strings.TrimSpace(ansi.Strip(stderr.String()))
		if message != "" {
			return nil, fmt.Errorf("list containers: %s: %w", message, err)
		}
		return nil, fmt.Errorf("list containers: %w", err)
	}

	var runtimeContainers []specs.Container
	if err := json.Unmarshal(output, &runtimeContainers); err != nil {
		return nil, fmt.Errorf("decode container list: %w", err)
	}

	return &specs.Containers{Items: runtimeContainers}, nil
}

// Selects project containers unless every runtime container was requested.
func selectProjectContainers(containers *specs.Containers, graph *specs.Graph, all bool) *specs.Containers {
	if all {
		return containers
	}

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

	selectedContainers := make([]specs.Container, 0, len(containerNames))
	for _, container := range containers.Items {
		if slices.Contains(containerNames, container.ID) {
			selectedContainers = append(selectedContainers, container)
		}
	}

	return &specs.Containers{Items: selectedContainers}
}
