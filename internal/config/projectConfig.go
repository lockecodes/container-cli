package config

// ProjectConfig defines the configuration for a specific project within the container CLI system.
// It includes details such as project name, file paths, and default settings for build and runtime.
type ProjectConfig struct {
	Name           string `koanf:"name"`
	Path           string `koanf:"path"`
	Dockerfile     string `koanf:"dockerfile"`
	BuildDirectory string `koanf:"buildDirectory"`
	BuildContext   string `koanf:"buildContext"`
	DefaultCommand string `koanf:"defaultCommand"`
	CommandAlias   string `koanf:"commandAlias"`
}
