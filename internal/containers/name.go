package containers

import "fmt"

// Builds the project-scoped FQDN used by Apple Container.
func serviceContainerName(projectName string, serviceName string, replica int) string {
	return fmt.Sprintf("%s-%d.%s", serviceName, replica, projectName)
}
