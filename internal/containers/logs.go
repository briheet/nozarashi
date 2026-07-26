package containers

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/render"
	"github.com/briheet/nozarashi/internal/specs"
)

// This function wraps over apple's container cli and helps in getting logs
func LogsContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running
	if err := StatusSystemContainers(ctx, io.Discard); err != nil {
		return err
	}

	// Parse the config file
	specs, err := config.ParseTOMLConfig(ctx, opts.FilePath)
	if err != nil {
		return err
	}

	// Validate dependency graph and find any inconsistencies
	if err := validateDependencyGraph(specs); err != nil {
		return err
	}

	// Build dependency graph
	graph, err := buildDependencyGraph(specs)
	if err != nil {
		return err
	}

	// Select services passed through the command arguments
	if err := selectServiceNodes(graph, opts.Args); err != nil {
		return err
	}

	// Get logs
	logs, err := getContainerLogs(ctx, graph, opts.Number)
	if err != nil {
		return err
	}

	// Render logs from every selected container.
	if err := render.RenderLogs(ctx, logs); err != nil {
		return fmt.Errorf("render container logs: %w", err)
	}

	return nil
}

// Gets cumulative logs from selected service containers.
func getContainerLogs(ctx context.Context, graph *specs.Graph, number int) (*specs.Logs, error) {
	logs := &specs.Logs{
		Containers: make([]specs.ContainerLogs, 0, len(graph.Nodes)),
	}

	for _, serviceNode := range graph.Nodes {
		replicas := max(serviceNode.Spec.Replicas, 1)

		for replica := 1; replica <= replicas; replica++ {
			containerName := serviceContainerName(
				graph.Project.Project.Name,
				serviceNode.ServiceName,
				replica,
			)

			lines, err := GetContainerLogLines(ctx, containerName, number)
			if err != nil {
				return nil, err
			}

			logs.Containers = append(logs.Containers, specs.ContainerLogs{
				Name:  containerName,
				Lines: lines,
			})
		}
	}

	return logs, nil
}

// GetContainerLogLines gets bounded logs from one container.
func GetContainerLogLines(ctx context.Context, name string, number int) ([]string, error) {
	logsCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerLogsArgs(name, number)...,
	)

	output, err := logsCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"get logs for container %q: %s: %w",
			name,
			strings.TrimSpace(string(output)),
			err,
		)
	}

	if logOutput := strings.TrimSuffix(string(output), "\n"); logOutput != "" {
		return strings.Split(logOutput, "\n"), nil
	}

	return nil, nil
}
