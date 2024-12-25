package container

import (
	"bytes"
	"fmt"
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"log"
	"os"
	"os/exec"
	"strings"
)

//const PODMAN_URI = "unix:///run/podman/podman.sock"

//const PODMAN_URI = "unix:/run/user/1000/podman/podman.sock"

type Podman struct {
	BuildContext              string
	BuildDirectory            string
	ContextDirectoryContainer string
	ContextDirectoryHost      string
	Dockerfile                string
	ImageName                 string
	ImageTag                  string
	UserHomeContainer         string
	UserHomeHost              string
	DefaultCommand            string
}

func (p Podman) Build() error {
	// build using shell for now
	// come back and use go for this later...maybe...
	err := buildImageUsingShell(p.Dockerfile, p.ImageName, p.BuildContext)
	if err != nil {
		return err
	}
	return nil
}

func (p Podman) Stop() error {
	panic("implement me")
}

// GetRunCommand returns the command used for running the application
func (p Podman) GetRunCommand() []string {
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
	commandStr := fmt.Sprintf("podman %s", strings.Join(cmdArgs, " "))
	fmt.Println("Generated Command:", commandStr)

	return cmdArgs
}

// Run generates and executes a `podman run` command with the given environment variables and volumes.
func (p Podman) Run(args []string) error {
	// Print the full command for logging or debugging
	cmdArgs := p.GetRunCommand()
	cmdArgs = append(cmdArgs, args...)

	// Execute the command
	cmd := exec.Command("podman", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

func NewPodman(projectConfig *config.ProjectConfig) Container {
	workingDir, err := os.Getwd() // Returns the directory from where the program is invoked
	if err != nil {
		log.Fatal(err)
	}
	homeDir, _ := os.UserHomeDir()
	return Podman{
		BuildContext:              projectConfig.BuildContext,
		BuildDirectory:            projectConfig.BuildDirectory,
		ContextDirectoryContainer: globals.CONTEXT_DIRECTORY_CONTAINER,
		ContextDirectoryHost:      workingDir,
		Dockerfile:                projectConfig.Dockerfile,
		ImageName:                 projectConfig.Name,
		ImageTag:                  "latest",
		UserHomeContainer:         globals.USER_HOME_CONTAINER,
		UserHomeHost:              homeDir,
		DefaultCommand:            projectConfig.DefaultCommand,
	}
}

// buildImageUsingShell performs a Podman build by invoking the host shell.
func buildImageUsingShell(dockerfilePath, imageName, contextDir string) error {
	// Construct the command arguments for `podman build`
	cmdArgs := []string{
		"build",
		"-f", dockerfilePath, // Specify the Dockerfile path
		"-t", imageName, // Assign the tag to the image
		contextDir, // Set the build context (working directory)
	}

	// Create a new Podman build command
	cmd := exec.Command("podman", cmdArgs...)

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
