package container

type Container interface {
	Build() error
	Run(args []string) error
	Stop() error
}
