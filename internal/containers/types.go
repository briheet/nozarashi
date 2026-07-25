package containers

// Options used by container lifecycle commands.
type ContainerOptions struct {
	// Passing in the filepath of the config file
	FilePath string

	// Arguments such as container names
	Args []string

	// Number of line for logs to print
	Number int

	// List all containers instead of project containers
	All bool

	// Replica selects a service container replica
	Replica int

	// Interactive keeps standard input open for exec
	Interactive bool

	// TTY allocates a terminal for exec
	TTY bool
}
