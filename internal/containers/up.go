package containers

import (
	"context"

	"github.com/briheet/nozarashi/internal/config"
)

// Options when using the up command
type UpOptions struct {
	// Passing in the filepath of the config file
	FilePath string

	// Container names that one wants
	Containers []string
}

// This wraps over apple's container cli and helps us creating and managing containers
func UpContainers(ctx context.Context, opts UpOptions) error {
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
	if err := validateDependencyGraph(ctx, specs, opts); err != nil {
		return err
	}

	// Build dependency graph
	graph, err := buildDependencyGraph(ctx, specs)
	if err != nil {
		return err
	}

	// Build Service Images
	if err := buildServiceImages(ctx, graph); err != nil {
		return err
	}

	// Create Resources
	if err := createResources(ctx, graph); err != nil {
		return err
	}

	// Create Service Containers
	if err := createImages(ctx, graph); err != nil {
		return err
	}

	return nil
}
