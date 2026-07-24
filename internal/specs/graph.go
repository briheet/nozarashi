package specs

// Enum type for Resources.
// This includes switch on container resources
type ResourceKind string

const (
	ResourceVolume  ResourceKind = "volume"
	ResourceNetwork ResourceKind = "network"
	ResourceConfig  ResourceKind = "config"
	ResourceSecret  ResourceKind = "secret"
)

// This is the main struct for dependency graph resolving and building
// All inputs will be setup, services named formatted, configs and secrets passed
type Graph struct {
	// Root node. Will hold config, secrets, volumes, network, inputs projects access
	Project *Specs

	// Prestarting resources. This includes volumes, networks, configs
	Resources []*ResourceNode

	// Service nodes. Will hold all service with access to everything required to spin up
	Nodes []*ServiceNode
}

// ServiceNode represents a service after its configuration,
// resource references and dependencies have been resolved.
type ServiceNode struct {
	// Logical service name from the project specification
	ServiceName string

	// Runtime container name
	// Formatted as <project-name>-<service-name>-<replicas>
	ContainerName string

	// Reference to the original service specification.
	Spec *ServiceSpecs

	// Names of services that must start before this service.
	DependsOn []string
}

// ResourceNode represents a project resource that must be
// prepared before service containers are created.
type ResourceNode struct {
	// Name of the resource
	// Formatted as <project-name>-<resource-name>
	Name string

	// Different different resource kinds defined in enum types
	Kind ResourceKind

	// Pointers to be passed from Specs.
	// Make sure to carefully switch on them lmao
	Volume  *VolumeSpecs
	Network *NetworkSpecs
	Config  *ConfigSpecs
	Secret  *SecretSpecs
}
