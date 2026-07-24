package containers

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/briheet/nozarashi/internal/specs"
	"github.com/distribution/reference"
)

// This function helps in validateing dependency graph
// This includes check for validating links, services and resource matchups, etc.
// Few things it does:
// Validate every depends_on points to an existing service
// Validate service runnable sources (image, build, nix input)
// Detect dependency cycles
// Validate referenced volumes, network, config, secret
// Verifies dependencies via lockfile (lockfile is relative to config file)
func validateDependencyGraph(projectSpecs *specs.Specs) error {
	// Validate all services have correct referencces to other services and resources
	if err := validateServiceReferences(projectSpecs); err != nil {
		return err
	}

	// Validate directed cyclic dependencies via DependsOn
	if err := validateDependencyCycles(projectSpecs.Services); err != nil {
		return err
	}

	return nil
}

// This function does a dfs checking for cycling as direct due to dependsOn
// TODO: Please benchmark and rewrite this
func validateDependencyCycles(services map[string]specs.ServiceSpecs) error {
	const (
		unvisited uint8 = iota
		visiting
		visited
	)

	color := make(map[string]uint8, len(services))
	parent := make(map[string]string, len(services))

	var cycleStart string
	var cycleEnd string

	var dfs func(string) bool

	dfs = func(serviceName string) bool {
		color[serviceName] = visiting

		for _, dependency := range services[serviceName].DependsOn {
			switch color[dependency] {
			case unvisited:
				parent[dependency] = serviceName

				if dfs(dependency) {
					return true
				}

			case visiting:
				cycleEnd = serviceName
				cycleStart = dependency
				return true
			}
		}

		color[serviceName] = visited
		return false
	}

	for serviceName := range services {
		if color[serviceName] == unvisited && dfs(serviceName) {
			break
		}
	}

	if cycleStart == "" {
		return nil
	}

	cycle := []string{cycleStart}

	for current := cycleEnd; current != cycleStart; current = parent[current] {
		cycle = append(cycle, current)
	}

	cycle = append(cycle, cycleStart)
	slices.Reverse(cycle)

	return fmt.Errorf(
		"service dependency cycle detected: %s",
		strings.Join(cycle, " -> "),
	)
}

// This helps mapping services to other services on being present
// This also maps services to required resources
func validateServiceReferences(specs *specs.Specs) error {
	// Loop over and check via dependsOn
	for serviceName, serviceSpec := range specs.Services {
		for _, dependency := range serviceSpec.DependsOn {
			if _, ok := specs.Services[dependency]; !ok {
				return fmt.Errorf("Dependency missing, Service: %s, Dependecy: %s", serviceName, dependency)
			}
		}
	}

	// Validate reference on the basis of types
	for _, serviceSpec := range specs.Services {
		switch serviceSpec.Type {
		case "input":
			if _, ok := specs.Inputs[serviceSpec.Reference]; !ok {
				return fmt.Errorf("input %q is not declared", serviceSpec.Reference)
			}

		case "path":
			if !filepath.IsLocal(serviceSpec.Reference) {
				return fmt.Errorf("path reference %q must remain inside the project", serviceSpec.Reference)
			}

		case "oci":
			if _, err := reference.ParseNormalizedNamed(serviceSpec.Reference); err != nil {
				return fmt.Errorf("invalid OCI reference %q: %w", serviceSpec.Reference, err)
			}

		default:
			return fmt.Errorf("unsupported image source %q", serviceSpec.Type)

		}
	}
	return nil
}
