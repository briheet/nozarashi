package containers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/render"
	"github.com/briheet/nozarashi/internal/specs"
)

// This function wraps over apple's container cli and helps in getting logs
func LogsContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running
	if err := StatusSystemContainers(ctx); err != nil {
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
			containerName := fmt.Sprintf("%s-%d", serviceNode.ServiceName, replica)

			logsCmd := exec.CommandContext(
				ctx,
				ContainerCliName,
				ContainerLogsArgs(containerName, number)...,
			)
			logsCmd.Stderr = os.Stderr

			output, err := logsCmd.Output()
			if err != nil {
				return nil, fmt.Errorf("get logs for container %q: %w", containerName, err)
			}

			lines := make([]string, 0)
			if logOutput := strings.TrimSuffix(string(output), "\n"); logOutput != "" {
				lines = strings.Split(logOutput, "\n")
			}

			logs.Containers = append(logs.Containers, specs.ContainerLogs{
				Name:  containerName,
				Lines: lines,
			})
		}
	}

	return logs, nil
}
