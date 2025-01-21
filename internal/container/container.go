package container

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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

// Build builds a Podman image using the specified Dockerfile, image name, and build context.
func (p Container) Build() error {
	// build using shell for now
	// come back and use go for this later...maybe...
	err := buildImageUsingShell(p.ContainerEngine, p.Dockerfile, p.ImageName, p.BuildContext)
	if err != nil {
		return err
	}
	return nil
}

// buildImageUsingShell performs a Podman build by invoking the host shell.
func buildImageUsingShell(engine, dockerfilePath, imageName, contextDir string) error {
	// Construct the command arguments for `podman build`
	cmdArgs := []string{
		"build",
		"-f", dockerfilePath, // Specify the Dockerfile path
		"-t", imageName, // Assign the tag to the image
		contextDir, // Set the build context (working directory)
	}

	// Create a new Podman build command
	cmd := exec.Command(engine, cmdArgs...)

	// Capture the standard output and standard error
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set the working directory to the context directory
	cmd.Dir = contextDir

	// Execute the Podman command
	log.Printf("Running command: podman %v\n", cmdArgs)
	err := cmd.Run()

	// Print any standard output and error
	if stdout.Len() > 0 {
		fmt.Printf("Output:\n%s\n", stdout.String())
	}
	if stderr.Len() > 0 {
		fmt.Printf("Error Output:\n%s\n", stderr.String())
	}

	// Check for errors in running the command
	if err != nil {
		return fmt.Errorf("failed to execute podman build: %w", err)
	}

	log.Println("Podman image built successfully!")
	return nil
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
