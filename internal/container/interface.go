package container

// Container defines an interface for managing container lifecycle operations such as build, run, and stop.
// Build constructs the container image from a specified configuration or context.
// Run executes the container with the provided arguments.
// GetRunCommand retrieves the command used to start the container.
// Stop halts a running container and cleans up associated resources.
type Container interface {
	Build() error
	Run(args []string) error
	GetRunCommand() []string
	Stop() error
}
