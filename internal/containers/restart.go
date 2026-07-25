package containers

import (
	"context"

	"github.com/briheet/nozarashi/internal/config"
)

// This function restarts project service containers.
func RestartContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running.
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

	// Select services passed through the command arguments.
	if err := selectServiceNodes(graph, opts.Args); err != nil {
		return err
	}

	// Stop selected containers in reverse dependency order.
	if err := stopContainers(ctx, graph); err != nil {
		return err
	}

	// Start selected containers in dependency order.
	if err := createContainers(ctx, graph); err != nil {
		return err
	}

	return nil
}
