package containers

import (
	"context"
	"fmt"
	"io"

	"github.com/briheet/nozarashi/internal/config"
	"golang.org/x/sync/errgroup"
)

// This function destroys containers, networks, volumes and images for a project.
func DestroyContainers(ctx context.Context, opts ContainerOptions) error {
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

	// Stop and delete service containers in reverse dependency order.
	if err := stopContainers(ctx, graph); err != nil {
		return err
	}

	// Remove containers
	if err := removeContainers(ctx, graph); err != nil {
		return err
	}

	// Delete independent project resources concurrently.
	group, groupCtx := errgroup.WithContext(ctx)

	// Delete project runtime networks.
	group.Go(func() error {
		return removeNetworks(groupCtx, graph)
	})

	// Delete project runtime volumes.
	group.Go(func() error {
		return removeVolumes(groupCtx, graph)
	})

	// Delete local and registry-backed service images.
	group.Go(func() error {
		return destroyImages(groupCtx, graph)
	})

	if err := group.Wait(); err != nil {
		return fmt.Errorf("destroy project resources: %w", err)
	}

	return nil
}
