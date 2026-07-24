package containers

import (
	"context"
	"fmt"
	"slices"

	"github.com/briheet/nozarashi/internal/specs"
)

// Build dependency graph
// As all resource are validated, it builds inline order dependencyGraph to execute
// We do a topological sort here. Check this for more: https://cp-algorithms.com/graph/topological-sort.html
// Graphs are already validated and directed (dependsOn), hence no contradiction
func buildDependencyGraph(ctx context.Context, projectSpecs *specs.Specs) (*specs.Graph, error) {
	// Base graph definition
	var graph specs.Graph

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
			// <project-name>-<driver>
			Name:   fmt.Sprintf("%s-%s", name, resource.Driver),
			Kind:   specs.ResourceVolume,
			Volume: &resource,
		}
		resourceNodes = append(resourceNodes, node)
	}

	// Second build up for networks
	for name, resource := range projectSpecs.Networks {
		node := &specs.ResourceNode{
			Name:    fmt.Sprintf("%s-%s", name, resource.Driver),
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

// Build Service images. If exists, then recreate
func buildServiceImages(ctx context.Context, graph *specs.Graph) error {
	return nil
}
