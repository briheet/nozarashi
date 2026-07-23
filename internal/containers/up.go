package containers

import "context"

// Options when using the up command
type UpOptions struct {
	// Passing in the filepath of the config file
	FilePath string

	// Container names that one wants
	Containers []string
}

// This wraps over apple's container cli and helps us creating and managing containers
func UpContainers(ctx context.Context, opts UpOptions) error {
	// First check this containers system is running
	if err := StatusSystemContainers(ctx); err != nil {
		return err
	}

	// Parse the config file

	return nil
}
