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
	// CreationDate is the time at which the container was created.
	CreationDate string `json:"creationDate"`

	// Image identifies the image used to create the container.
	Image ContainerImage `json:"image"`

	// Platform identifies the configured operating system and architecture.
	Platform ContainerPlatform `json:"platform"`

	// Resources contains the configured CPU and memory limits.
	Resources ContainerResources `json:"resources"`

	// InitProcess contains the container's initial process configuration.
	InitProcess ContainerInitProcess `json:"initProcess"`

	// Mounts contains the filesystems attached to the container.
	Mounts []ContainerMount `json:"mounts"`

	// Networks contains the configured network attachments.
	Networks []ContainerNetwork `json:"networks"`

	// PublishedPorts contains host-to-container port mappings.
	PublishedPorts []ContainerPublishedPort `json:"publishedPorts"`

	// Labels contains metadata attached to the container.
	Labels map[string]string `json:"labels"`
}

// ContainerImage identifies the image used by a container.
type ContainerImage struct {
	// Reference is the image name and optional tag or digest.
	Reference string `json:"reference"`
}

// ContainerPlatform identifies the container operating system and architecture.
type ContainerPlatform struct {
	// OperatingSystem is the operating system expected by the image.
	OperatingSystem string `json:"os"`

	// Architecture is the CPU architecture expected by the image.
	Architecture string `json:"architecture"`
}

// ContainerResources contains allocated compute resources.
type ContainerResources struct {
	// CPUs is the number of virtual CPUs allocated to the container.
	CPUs int `json:"cpus"`

	// MemoryInBytes is the container's configured memory limit.
	MemoryInBytes int64 `json:"memoryInBytes"`
}

// ContainerInitProcess contains the configured process and environment.
type ContainerInitProcess struct {
	// Executable is the program started inside the container.
	Executable string `json:"executable"`

	// Arguments contains arguments passed to the executable.
	Arguments []string `json:"arguments"`

	// Environment contains the process environment as KEY=value entries.
	Environment []string `json:"environment"`

	// WorkingDirectory is the process's initial directory.
	WorkingDirectory string `json:"workingDirectory"`
}

// ContainerMount contains a container filesystem mount.
type ContainerMount struct {
	// Source is the host path used by a bind mount.
	Source string `json:"source"`

	// Destination is the path at which the mount appears in the container.
	Destination string `json:"destination"`

	// Type contains type-specific mount information.
	Type ContainerMountType `json:"type"`
}

// ContainerMountType contains named volume information.
type ContainerMountType struct {
	// Volume identifies a named volume mount.
	Volume ContainerVolumeMount `json:"volume"`
}

// ContainerVolumeMount identifies a named container volume.
type ContainerVolumeMount struct {
	// Name is the managed volume name.
	Name string `json:"name"`
}

// ContainerNetwork contains a configured network attachment.
type ContainerNetwork struct {
	// Network is the managed network name.
	Network string `json:"network"`
}

// ContainerPublishedPort contains a host-to-container port mapping.
type ContainerPublishedPort struct {
	// HostAddress is the host interface on which the port is published.
	HostAddress string `json:"hostAddress"`

	// HostPort is the port exposed on the host.
	HostPort int `json:"hostPort"`

	// ContainerPort is the destination port inside the container.
	ContainerPort int `json:"containerPort"`

	// Protocol is the transport protocol used by the mapping.
	Protocol string `json:"proto"`
}

// ContainerStatus contains current state and network information.
type ContainerStatus struct {
	// State is the current container lifecycle state.
	State string `json:"state"`

	// StartedDate is the time at which the container last started.
	StartedDate string `json:"startedDate"`

	// Networks contains runtime addresses assigned to the container.
	Networks []ContainerNetworkStatus `json:"networks"`
}

// ContainerNetworkStatus contains a runtime network address.
type ContainerNetworkStatus struct {
	// Network is the managed network name.
	Network string `json:"network"`

	// IPv4Address is the assigned IPv4 address and prefix.
	IPv4Address string `json:"ipv4Address"`

	// IPv6Address is the assigned IPv6 address and prefix.
	IPv6Address string `json:"ipv6Address"`

	// MacAddress is the interface's assigned hardware address.
	MacAddress string `json:"macAddress"`
}

// ContainerStats contains live resource usage counters.
type ContainerStats struct {
	// ID is the runtime container name.
	ID string `json:"id"`

	// CPUUsageUsec is cumulative CPU time consumed by the container.
	CPUUsageUsec uint64 `json:"cpuUsageUsec"`

	// MemoryUsageBytes is the container's current memory usage.
	MemoryUsageBytes uint64 `json:"memoryUsageBytes"`

	// MemoryLimitBytes is the container's configured memory limit.
	MemoryLimitBytes uint64 `json:"memoryLimitBytes"`

	// NetworkRxBytes is the cumulative number of network bytes received.
	NetworkRxBytes uint64 `json:"networkRxBytes"`

	// NetworkTxBytes is the cumulative number of network bytes transmitted.
	NetworkTxBytes uint64 `json:"networkTxBytes"`

	// BlockReadBytes is the cumulative number of bytes read from block devices.
	BlockReadBytes uint64 `json:"blockReadBytes"`

	// BlockWriteBytes is the cumulative number of bytes written to block devices.
	BlockWriteBytes uint64 `json:"blockWriteBytes"`

	// NumProcesses is the current number of processes inside the container.
	NumProcesses int `json:"numProcesses"`

	// CPUPercent is CPU usage derived from two consecutive samples.
	CPUPercent float64 `json:"-"`

	// CPUPercentValid reports whether CPUPercent has a previous sample.
	CPUPercentValid bool `json:"-"`
}
