package specs

// Enum type for Inputs
type InputType string

const (
	InputTypeGit   InputType = "git"
	InputTypeNix   InputType = "flake"
	InputTypeLocal InputType = "path"
)

// Enum type for Services
type ServiceType string

const (
	ServiceTypeOCI   ServiceType = "oci"
	ServiceTypeInput ServiceType = "input"
	ServiceTypePath  ServiceType = "path"
)

// This will contain standard specs.
// For more info please check docs/SPEC.md
type Specs struct {
	// Project Metadata and defaults.
	Project ProjectSpecs `mapstructure:"project" validate:"required"`

	// Pinned external Nix, Git or module inputs.
	Inputs map[string]InputSpecs `mapstructure:"inputs" validate:"required"`

	// Long-running project services.
	Services map[string]ServiceSpecs `mapstructure:"services" validate:"required"`

	// Named persistent volumes.
	Volumes map[string]VolumeSpecs `mapstructure:"volumes" validate:"required"`

	// User-defined networks.
	Networks map[string]NetworkSpecs `mapstructure:"networks" validate:"required"`

	// Configs need to be provided to services
	Configs map[string]ConfigSpecs `mapstructure:"configs"`

	// Secrets that need to be provided to services
	Secrets map[string]SecretSpecs `mapstructure:"secrets"`
}

// This defines project's types which can be defined
type ProjectSpecs struct {
	// Name of the project
	Name string `mapstructure:"name" validate:"required"`

	// Version of the project
	Version string `mapstructure:"version" validate:"required"`

	// Description of the project
	Description string `mapstructure:"description" validate:"required"`

	// Profile of the project
	// This could be anything from dev to prod to main
	// This is added so that if one shares, it can have multiple files and get its type
	Profile string `mapstructure:"profile" validate:"required"`
}

// This defines input specs types
type InputSpecs struct {
	// Type of input whether it is a flake, Git repository or local Nix module
	Type InputType `mapstructure:"type" validate:"required"`

	// Source of the input type. Can be Url or a relative path
	Source string `mapstructure:"source"`

	// Ref is a reference for Nix, Git
	// This includes branch, commit, tag
	Ref string `mapstructure:"ref"`
}

// This defines the service specs types
type ServiceSpecs struct {
	// Type of service.
	// oci, path, input
	Type ServiceType `mapstructure:"type" validate:"required,oneof=oci path input"`

	// Reference to type
	Reference string `mapstructure:"reference" validate:"required"`

	// Attribute if so
	Attribute string `mapstructure:"attribute"`

	// Entrypoint to the application
	EntryPoint []string `mapstructure:"entrypoint"`

	// Command to execute
	Command []string `mapstructure:"command"`

	// Environment variables passed on to the env
	Environment map[string]string `mapstructure:"environment"`

	// Persistent volumes to the service
	Volumes []VolumeMount `mapstructure:"volumes"`

	// Networks this service joins
	Networks []string `mapstructure:"networks"`

	// Exposed ports
	Ports []PortSpec `mapstructure:"ports"`

	// Service dependencies
	DependsOn []string `mapstructure:"depends_on"`

	// Number of replicas
	Replicas int `mapstructure:"replicas" validate:"omitempty,min=1,max=100"`

	// Policy
	Policy string `mapstructure:"policy" validate:"omitempty,oneof=no always on-failure unless-stopped"`
}

// VolumeSpecs defines a named persistent volume.
type VolumeSpecs struct {
	// Driver controls how the volume is created.
	// Initially, this can support "local".
	Driver string `mapstructure:"driver"`
}

// NetworkSpecs defines a named container network.
type NetworkSpecs struct {
	// Driver determines the network implementation.
	// Examples: "bridge" or "host".
	Driver string `mapstructure:"driver"`
}

// ConfigSpecs defines non-sensitive configuration supplied to services.
type ConfigSpecs struct {
	// Name of the config object to look into
	Name string `mapstructure:"name"`

	// File is a local file containing the configuration.
	File string `mapstructure:"file"`

	// Content can be created with inline values
	Content string `mapstructure:"content"`

	// External means the secret is managed outside
	External bool `mapstructure:"external"`
}

// SecretSpecs defines sensitive data supplied to services.
type SecretSpecs struct {
	// Name of the config object to look into
	Name string `mapstructure:"name"`

	// File is a local file containing the configuration
	File string `mapstructure:"file"`

	// External means the secret is managed outside
	External bool `mapstructure:"external"`
}

// Volume mount for Service
type VolumeMount struct {
	// Name of the volume mount
	Name string `mapstructure:"name"`

	// Path to store
	Path string `mapstructure:"path"`
}

// Port specs for the service
type PortSpec struct {
	// Host on
	Host int `mapstructure:"host"`

	// Container specific
	Container int `mapstructure:"container"`
}
