package globals

import (
	"os"
	"path/filepath"
)

var HomeDir string
var DefaultContainerCliConfigPath string

const ProjectId string = "47137983"

func init() {
	HomeDir, _ = os.UserHomeDir()
	DefaultContainerCliConfigPath = filepath.Join(HomeDir, ".config/container-cli/config.yaml")
}

const USER_HOME_CONTAINER = "/opt/usr/home"
const CONTEXT_DIRECTORY_CONTAINER = "/opt/context"
