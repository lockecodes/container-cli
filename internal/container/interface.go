package container

type Container interface {
	Build() error
	Run(args []string) error
	GetRunCommand() []string
	Stop() error
}
