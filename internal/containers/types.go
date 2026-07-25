package containers

// Options used by container lifecycle commands.
type ContainerOptions struct {
	// Passing in the filepath of the config file
	FilePath string

	// Arguments such as container names
	Args []string

	// Number of line for logs to print
	Number int
}
