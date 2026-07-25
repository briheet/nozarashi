package containers

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/briheet/nozarashi/internal/config"
	"github.com/briheet/nozarashi/internal/specs"
)

// This function wraps over apple's container cli and executes a command in a service.
func ExecContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running.
	if err := StatusSystemContainers(ctx); err != nil {
		return err
	}

	// Integration tests check
	if len(opts.Args) < 2 {
		return fmt.Errorf("exec requires a service and command")
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

	// Select the service passed as the first positional argument.
	serviceName := opts.Args[0]
	if err := selectServiceNodes(graph, []string{serviceName}); err != nil {
		return err
	}

	return execServiceContainer(
		ctx,
		graph.Project.Project.Name,
		serviceName,
		graph.Nodes[0],
		opts,
	)
}

// Executes a command in the selected service replica.
func execServiceContainer(
	ctx context.Context,
	projectName string,
	serviceName string,
	serviceNode *specs.ServiceNode,
	opts ContainerOptions,
) error {
	replicas := max(serviceNode.Spec.Replicas, 1)
	if opts.Replica < 1 || opts.Replica > replicas {
		return fmt.Errorf(
			"replica %d for service %q must be between 1 and %d",
			opts.Replica,
			serviceName,
			replicas,
		)
	}

	containerName := serviceContainerName(
		projectName,
		serviceNode.ServiceName,
		opts.Replica,
	)
	execCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerExecArgs(containerName, opts)...,
	)

	if opts.Interactive {
		execCmd.Stdin = os.Stdin
	}
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("execute command in container %q: %w", containerName, err)
	}

	return nil
}
