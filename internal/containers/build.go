package containers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/briheet/nozarashi/internal/nix"
	"github.com/briheet/nozarashi/internal/specs"
)

// Build dependency graph
// As all resource are validated, it builds inline order dependencyGraph to execute
// We do a topological sort here. Check this for more: https://cp-algorithms.com/graph/topological-sort.html
// Graphs are already validated and directed (dependsOn), hence no contradiction
func buildDependencyGraph(ctx context.Context, projectSpecs *specs.Specs) (*specs.Graph, error) {
	// Base graph definition
	var graph specs.Graph

	// Copy the project specs
	graph.Project = projectSpecs

	// Build resource nodes and save
	graph.Resources = buildResourceNodes(projectSpecs)

	// Build Service nodes
	graph.Nodes = buildServiceNodes(projectSpecs)

	return &graph, nil
}

// Builds resource nodes and sends back
// Currently only handles Volumes and Networks
func buildResourceNodes(projectSpecs *specs.Specs) []*specs.ResourceNode {
	resourceNodes := make([]*specs.ResourceNode, 0, len(projectSpecs.Volumes)+len(projectSpecs.Networks))

	// First buildup for volumes
	for name, resource := range projectSpecs.Volumes {
		node := &specs.ResourceNode{
			// <project-name>-<resource-name>
			Name:   fmt.Sprintf("%s-%s", projectSpecs.Project.Name, name),
			Kind:   specs.ResourceVolume,
			Volume: &resource,
		}
		resourceNodes = append(resourceNodes, node)
	}

	// Second build up for networks
	for name, resource := range projectSpecs.Networks {
		node := &specs.ResourceNode{
			Name:    fmt.Sprintf("%s-%s", projectSpecs.Project.Name, name),
			Kind:    specs.ResourceNetwork,
			Network: &resource,
		}
		resourceNodes = append(resourceNodes, node)
	}

	// TODO: Config and secrets in future
	return resourceNodes
}

// Build nodes and topological sort
// TODO: Please rewrite this
func buildServiceNodes(projectSpecs *specs.Specs) []*specs.ServiceNode {
	serviceNodes := make([]*specs.ServiceNode, 0, len(projectSpecs.Services))

	visited := make(map[string]bool, len(projectSpecs.Services))

	var visit func(string)

	visit = func(serviceName string) {
		if visited[serviceName] {
			return
		}

		visited[serviceName] = true
		service := projectSpecs.Services[serviceName]

		dependencies := slices.Clone(service.DependsOn)
		slices.Sort(dependencies)

		for _, dependency := range dependencies {
			visit(dependency)
		}

		serviceSpec := service
		serviceNodes = append(serviceNodes, &specs.ServiceNode{
			ServiceName: fmt.Sprintf("%s-%s", projectSpecs.Project.Name, serviceName),
			Spec:        &serviceSpec,
		})
	}

	serviceNames := make([]string, 0, len(projectSpecs.Services))

	for serviceName := range projectSpecs.Services {
		serviceNames = append(serviceNames, serviceName)
	}

	slices.Sort(serviceNames)

	for _, serviceName := range serviceNames {
		visit(serviceName)
	}

	return serviceNodes
}

// Builds or pulls all service images.
func buildServiceImages(ctx context.Context, graph *specs.Graph) error {
	// Build or pull every service image in dependency order.
	for _, serviceNode := range graph.Nodes {
		switch serviceNode.Spec.Type {
		case specs.ServiceTypeInput:
			if err := buildNixServiceImage(ctx, graph, serviceNode); err != nil {
				return err
			}
		case specs.ServiceTypeOCI:
			if err := pullOCIServiceImage(ctx, serviceNode); err != nil {
				return err
			}
		case specs.ServiceTypePath:
			if err := buildContainerfileServiceImage(ctx, serviceNode); err != nil {
				return err
			}
		}
	}

	return nil
}

// Builds a Nix-backed service image and loads it into Apple Containers.
func buildNixServiceImage(ctx context.Context, graph *specs.Graph, serviceNode *specs.ServiceNode) error {
	// Build the image archive through dockerTools.buildLayeredImage + skeopo
	input := graph.Project.Inputs[serviceNode.Spec.Reference]
	archivePath, err := nix.BuildServiceImage(ctx, serviceNode.ServiceName, input, serviceNode.Spec)
	if err != nil {
		return fmt.Errorf("build Nix service image %q: %w", serviceNode.ServiceName, err)
	}

	// Load the built archive into the local image store.
	imageLoadCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerImageLoadArgs(archivePath)...,
	)
	imageLoadCmd.Stderr = os.Stderr
	imageLoadCmd.Stdout = os.Stdout

	if err := imageLoadCmd.Run(); err != nil {
		return fmt.Errorf("load Nix service image %q: %w", serviceNode.ServiceName, err)
	}

	return nil
}

// Pulls an OCI service image into the local image store.
func pullOCIServiceImage(ctx context.Context, serviceNode *specs.ServiceNode) error {
	imagePullCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerImagePullArgs(
			serviceNode.Spec.Reference,
			fmt.Sprintf(
				"%s/%s",
				specs.CurrentPlatform.OperatingSystem,
				specs.CurrentPlatform.Architecture,
			),
		)...,
	)
	imagePullCmd.Stderr = os.Stderr
	imagePullCmd.Stdout = os.Stdout

	if err := imagePullCmd.Run(); err != nil {
		return fmt.Errorf("pull OCI service image %q: %w", serviceNode.ServiceName, err)
	}

	return nil
}

// Builds a service image from a Dockerfile or Containerfile context.
func buildContainerfileServiceImage(ctx context.Context, serviceNode *specs.ServiceNode) error {
	// Use Dockerfile when present and fall back to Containerfile.
	filePath := filepath.Join(serviceNode.Spec.Reference, "Dockerfile")
	if _, err := os.Stat(filePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read Dockerfile for service %q: %w", serviceNode.ServiceName, err)
		}
		filePath = filepath.Join(serviceNode.Spec.Reference, "Containerfile")
	}

	imageBuildCmd := exec.CommandContext(
		ctx,
		ContainerCliName,
		ContainerImageBuildArgs(serviceNode.ServiceName, filePath, serviceNode.Spec.Reference)...,
	)
	imageBuildCmd.Stderr = os.Stderr
	imageBuildCmd.Stdout = os.Stdout

	if err := imageBuildCmd.Run(); err != nil {
		return fmt.Errorf("build Dockerfile or Containerfile service image %q: %w", serviceNode.ServiceName, err)
	}

	return nil
}
