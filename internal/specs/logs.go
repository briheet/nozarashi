package specs

// Logs contains collected stdio logs for service containers.
type Logs struct {
	// Containers keeps logs in service and replica order.
	Containers []ContainerLogs
}

// ContainerLogs contains line-oriented logs for one container.
type ContainerLogs struct {
	// Name is the project-scoped runtime container name.
	Name string

	// Lines contains container stdio split by line.
	Lines []string
}
