package containers

import (
	"context"

	"github.com/briheet/nozarashi/internal/config"
)

// This function destroys containers, networks, volumes and images for a project.
func DestroyContainers(ctx context.Context, opts ContainerOptions) error {
	// First check this containers system is running.
	if err := StatusSystemContainers(ctx); err != nil {
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

	// Stop and delete service containers in reverse dependency order.
	if err := stopContainers(ctx, graph); err != nil {
		return err
	}

	// Remove containers
	if err := removeContainers(ctx, graph); err != nil {
		return err
	}

	// Delete project runtime networks
	if err := removeNetworks(ctx, graph); err != nil {
		return err
	}

	// Delete project runtime volumes
	if err := removeVolumes(ctx, graph); err != nil {
		return err
	}

	// Delete local and registry-backed service images last.
	if err := destroyImages(ctx, graph); err != nil {
		return err
	}

	return nil
}
