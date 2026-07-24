package containers

// Options used by container lifecycle commands.
type ContainerOptions struct {
	// Passing in the filepath of the config file
	FilePath string

	// Rebuilding images Switch.
	// We should not always build images until passed as flag
	Build bool

	// Arguments such as container names
	Args []string
}
