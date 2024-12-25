package globals

import (
	"os"
	"path/filepath"
)

// HomeDir represents the user's home directory path, typically initialized using os.UserHomeDir().
var HomeDir string

// DefaultContainerCliConfigPath represents the default path for the container CLI configuration file.
var DefaultContainerCliConfigPath string

// ProjectId is a constant representing the unique identifier for the project in the container CLI configuration.
const ProjectId string = "47137983"

// UserHomeContainer defines the directory path in the container where the user's home directory is mapped.
const UserHomeContainer = "/opt/usr/home"

// ContextDirectoryContainer defines the directory path in the container where the runtime context is mounted.
const ContextDirectoryContainer = "/opt/context"

// init initializes the HomeDir and DefaultContainerCliConfigPath variables with appropriate default values.
func init() {
	HomeDir, _ = os.UserHomeDir()
	DefaultContainerCliConfigPath = filepath.Join(HomeDir, ".config/container-cli/config.yaml")
}
