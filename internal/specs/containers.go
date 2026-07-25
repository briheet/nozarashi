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
	CreationDate   string                   `json:"creationDate"`
	Image          ContainerImage           `json:"image"`
	Platform       ContainerPlatform        `json:"platform"`
	Resources      ContainerResources       `json:"resources"`
	InitProcess    ContainerInitProcess     `json:"initProcess"`
	Mounts         []ContainerMount         `json:"mounts"`
	Networks       []ContainerNetwork       `json:"networks"`
	PublishedPorts []ContainerPublishedPort `json:"publishedPorts"`
	Labels         map[string]string        `json:"labels"`
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

// ContainerInitProcess contains the configured process and environment.
type ContainerInitProcess struct {
	Executable       string   `json:"executable"`
	Arguments        []string `json:"arguments"`
	Environment      []string `json:"environment"`
	WorkingDirectory string   `json:"workingDirectory"`
}

// ContainerMount contains a container filesystem mount.
type ContainerMount struct {
	Source      string             `json:"source"`
	Destination string             `json:"destination"`
	Type        ContainerMountType `json:"type"`
}

// ContainerMountType contains named volume information.
type ContainerMountType struct {
	Volume ContainerVolumeMount `json:"volume"`
}

// ContainerVolumeMount identifies a named container volume.
type ContainerVolumeMount struct {
	Name string `json:"name"`
}

// ContainerNetwork contains a configured network attachment.
type ContainerNetwork struct {
	Network string `json:"network"`
}

// ContainerPublishedPort contains a host-to-container port mapping.
type ContainerPublishedPort struct {
	HostAddress   string `json:"hostAddress"`
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"proto"`
}

// ContainerStatus contains current state and network information.
type ContainerStatus struct {
	State       string                   `json:"state"`
	StartedDate string                   `json:"startedDate"`
	Networks    []ContainerNetworkStatus `json:"networks"`
}

// ContainerNetworkStatus contains a runtime network address.
type ContainerNetworkStatus struct {
	Network     string `json:"network"`
	IPv4Address string `json:"ipv4Address"`
	IPv6Address string `json:"ipv6Address"`
	MacAddress  string `json:"macAddress"`
}
