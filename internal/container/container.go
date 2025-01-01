package container

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/globals"
)

// Container represents configuration and attributes for managing container settings.
type Container struct {
	BuildContext              string
	BuildDirectory            string
	ContainerEngine           string
	ContextDirectoryContainer string
	ContextDirectoryHost      string
	Dockerfile                string
	ImageName                 string
	ImageTag                  string
	UserHomeContainer         string
	UserHomeHost              string
	DefaultCommand            string
}

// GetRunCommand returns the command used for running the application
func (p Container) GetRunCommand() []string {
	// Define environment variables
	envVars := map[string]string{
		"CONTEXT_DIR": p.ContextDirectoryContainer,
		"VERSION":     p.ImageTag,
		"IN_DOCKER":   "true",
		"DISPLAY":     os.Getenv("DISPLAY"),
	}

	// Define volume mappings
	volumes := []string{
		fmt.Sprintf("%s:%s", p.UserHomeHost, p.UserHomeContainer),                 // Map USER_HOME to /opt/usr/home
		fmt.Sprintf("%s:%s", p.ContextDirectoryHost, p.ContextDirectoryContainer), // Map CONTEXT_DIR to /opt/context
	}

	// Construct the `podman run` command
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "run") // Base command: "podman run"

	// Add environment variable flags
	for key, val := range envVars {
		if val != "" { // Ensure the variable is not empty
			cmdArgs = append(cmdArgs, "--env", fmt.Sprintf("%s=%s", key, val))
		}
	}

	// Add volume flags
	for _, volume := range volumes {
		if !strings.Contains(volume, ":") || strings.HasPrefix(volume, ":") {
			continue // Skip invalid volume mappings
		}
		cmdArgs = append(cmdArgs, "--volume", volume)
	}

	// Specify the image to run
	cmdArgs = append(cmdArgs, p.ImageName, p.DefaultCommand)

	// Print the full command for logging or debugging
	commandStr := fmt.Sprintf("%s %s", p.ContainerEngine, strings.Join(cmdArgs, " "))
	fmt.Println("Generated Command:", commandStr)

	return cmdArgs
}

// NewContainer initializes and returns a Container object configured with the given project settings.
func NewContainer(projectConfig *config.ProjectConfig) Container {
	workingDir, err := os.Getwd() // Returns the directory from where the program is invoked
	if err != nil {
		log.Fatal(err)
	}
	homeDir, _ := os.UserHomeDir()
	containerEngine, err := config.GetContainerEngineFromConfig()
	if err != nil {
		fmt.Println("failed to get container engine from config: ", err)
		panic(err)
	}
	return Container{
		BuildContext:              projectConfig.BuildContext,
		BuildDirectory:            projectConfig.BuildDirectory,
		ContainerEngine:           containerEngine,
		ContextDirectoryContainer: globals.ContextDirectoryContainer,
		ContextDirectoryHost:      workingDir,
		Dockerfile:                projectConfig.Dockerfile,
		ImageName:                 projectConfig.Name,
		ImageTag:                  "latest",
		UserHomeContainer:         globals.UserHomeContainer,
		UserHomeHost:              homeDir,
		DefaultCommand:            projectConfig.DefaultCommand,
	}
}
