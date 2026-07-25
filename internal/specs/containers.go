package specs

// Containers contains runtime containers for rendering.
type Containers struct {
	// Items keeps containers in runtime list order.
	Items []Container
}

// Container contains the Apple Container fields displayed by ps.
type Container struct {
	// ID is the runtime container name.
	ID string `json:"id"`

	// Configuration contains immutable container settings.
	Configuration ContainerConfiguration `json:"configuration"`

	// Status contains the current runtime state.
	Status ContainerStatus `json:"status"`
}

// ContainerConfiguration contains image, platform and resource settings.
type ContainerConfiguration struct {
	Image     ContainerImage     `json:"image"`
	Platform  ContainerPlatform  `json:"platform"`
	Resources ContainerResources `json:"resources"`
}

// ContainerImage identifies the image used by a container.
type ContainerImage struct {
	Reference string `json:"reference"`
}

// ContainerPlatform identifies the container operating system and architecture.
type ContainerPlatform struct {
	OperatingSystem string `json:"os"`
	Architecture    string `json:"architecture"`
}

// ContainerResources contains allocated compute resources.
type ContainerResources struct {
	CPUs          int   `json:"cpus"`
	MemoryInBytes int64 `json:"memoryInBytes"`
}

// ContainerStatus contains current state and network information.
type ContainerStatus struct {
	State       string                   `json:"state"`
	StartedDate string                   `json:"startedDate"`
	Networks    []ContainerNetworkStatus `json:"networks"`
}

// ContainerNetworkStatus contains a runtime network address.
type ContainerNetworkStatus struct {
	IPv4Address string `json:"ipv4Address"`
}
