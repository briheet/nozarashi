package containers

import (
	"fmt"
	"slices"

	"github.com/briheet/nozarashi/internal/specs"
)

// Selects service nodes passed as command arguments.
// An empty selection keeps every service in the graph.
func selectServiceNodes(graph *specs.Graph, serviceNames []string) error {
	// An empty argument list means every service remains selected.
	if len(serviceNames) == 0 {
		return nil
	}

	// Resolve logical service arguments to project-scoped runtime names.
	selectedServices := make([]string, 0, len(serviceNames))

	for _, serviceName := range serviceNames {
		// Reject unknown services before changing the graph.
		if _, ok := graph.Project.Services[serviceName]; !ok {
			return fmt.Errorf("service %q is not declared", serviceName)
		}

		selectedServices = append(
			selectedServices,
			fmt.Sprintf("%s-%s", graph.Project.Project.Name, serviceName),
		)
	}

	// Preserve dependency order while keeping only selected service nodes.
	serviceNodes := make([]*specs.ServiceNode, 0, len(selectedServices))
	for _, serviceNode := range graph.Nodes {
		if slices.Contains(selectedServices, serviceNode.ServiceName) {
			serviceNodes = append(serviceNodes, serviceNode)
		}
	}

	graph.Nodes = serviceNodes
	return nil
}
