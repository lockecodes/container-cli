package install

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/container"
	"gitlab.com/locke-codes/container-cli/internal/gitter"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/utils"
	"net/url"
	"os"
	"path"
	"strings"
)

// Project represents a project with a name, URL, destination directory, and a default command.
type Project struct {
	Name                 string
	URL                  string
	DestinationDirectory string
	DefaultCommand       string
}

// Path constructs and returns the full path of the project by combining the destination directory and project name.
func (p *Project) Path() string {
	return path.Join(p.DestinationDirectory, p.Name)
}

// Clone clones the project repository from the specified URL into the designated destination directory.
func (p *Project) Clone() error {
	fmt.Printf("Cloning %s\n", p.URL)
	client := gitter.NewGitter(p.Name, p.URL, p.Path())
	err := client.Clone()
	if err != nil {
		return err
	}
	return nil
}

// Install executes the installation process for the project, including cloning, configuring, and setting up scripts.
func (p *Project) Install() error {
	fmt.Printf("Installing %s\n", p.Name)
	err := p.Uninstall()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = p.Clone()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = p.InstallConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = p.InstallScript()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return err
}

// InstallScript creates and installs an executable script for the project in the user's local bin directory.
func (p *Project) InstallScript() error {
	projectConfig := config.ProjectConfig{
		Name:           p.Name,
		Path:           p.Path(),
		Dockerfile:     path.Join(p.Path(), "Dockerfile"),
		BuildDirectory: p.Path(),
		BuildContext:   p.Path(),
		DefaultCommand: p.DefaultCommand,
	}
	containerObj := container.NewPodman(&projectConfig)
	runCmd := containerObj.GetRunCommand()
	commandStr := fmt.Sprintf("podman %s", strings.Join(runCmd, " "))
	// File contents
	fileContent := fmt.Sprintf(`#!/usr/bin/env bash
%s $*`, commandStr)

	filePath := path.Join(globals.HomeDir, ".local/bin", p.Name)
	// Write the file content
	err = os.WriteFile(filePath, []byte(fileContent), 0755)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	// Make the file executable
	err = os.Chmod(filePath, 0755)
	if err != nil {
		fmt.Println("Error setting executable permissions:", err)
		return err
	}

	fmt.Printf("Script created and made executable at: %s\n", filePath)
	return nil
}

// InstallConfig installs or updates the project configuration in the container CLI configuration file.
func (p *Project) InstallConfig() error {
	fmt.Printf("Installing config for %s\n", p.Name)
	configFile := config.NewContainerCliConfig()
	err = configFile.LoadConfig()
	existingProject := configFile.GetProject(p.Name)
	projectConfig := config.ProjectConfig{
		Name:           p.Name,
		Path:           p.Path(),
		Dockerfile:     path.Join(p.Path(), "Dockerfile"),
		BuildDirectory: p.Path(),
		BuildContext:   p.Path(),
		DefaultCommand: p.DefaultCommand,
	}
	if existingProject == nil {
		configFile.Projects = append(configFile.Projects, projectConfig)
	} else {
		fmt.Printf("Project %s already exists\n", p.Name)
		err := configFile.ReplaceProjectByName(p.Name, projectConfig)
		if err != nil {
			return err
		}
	}
	configFile.KoanfLoad()
	err := configFile.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

// Uninstall removes all files and directories related to the project at the constructed project path.
// TODO: Also remove any symlinks and scripts
func (p *Project) Uninstall() error {
	fmt.Printf("Removing %s\n", p.Path())
	err := os.RemoveAll(p.Path())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return nil
}

// promptName prompts the user for a valid project name if not provided and validates the input for disallowed characters.
func promptName(name string) (string, error) {
	var err error

	validate := func(input string) error {
		chars := "!@#*+$&%\\/=~ \t\n"

		contains := strings.ContainsAny(input, chars)
		if contains {
			return fmt.Errorf("invalid name containers: %s", chars)
		}
		return nil
	}

	if name == "" {
		prompt := promptui.Prompt{
			Label:    "Name",
			Validate: validate,
		}

		name, err = prompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
	} else {
		return name, validate(name)
	}

	return name, nil
}

// promptCommand prompts the user for a valid command if not provided and validates the input for disallowed characters.
func promptCommand(command string) (string, error) {
	var err error

	validate := func(input string) error {
		// TODO: Should we make the command less restrictive for characters?
		chars := "!@#*+$&%\\/=~ \t\n"

		contains := strings.ContainsAny(input, chars)
		if contains {
			return fmt.Errorf("invalid command containers: %s", chars)
		}
		return nil
	}

	if command == "" {
		prompt := promptui.Prompt{
			Label:    "Command",
			Validate: validate,
		}

		command, err = prompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
	} else {
		return command, validate(command)
	}

	return command, nil
}

// promptUrl prompts the user for a valid project URL if not provided and validates the input for scheme and format.
func promptUrl(projectUrl string) (string, error) {
	var err error
	validate := func(input string) error {
		// Parse the input string as a URL
		parsedURL, err := url.Parse(input)
		if err != nil {
			return fmt.Errorf("invalid URL format: %w", err)
		}

		// Ensure the URL scheme is "https"
		if parsedURL.Scheme != "https" && parsedURL.Scheme != "ssh" {
			return fmt.Errorf("invalid URL scheme: %s (only https and ssh is allowed)", parsedURL.Scheme)
		}

		// No errors, the URL is valid
		return nil
	}
	if projectUrl == "" {
		prompt := promptui.Prompt{
			Label:    "URL",
			Validate: validate,
		}
		projectUrl, err = prompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
	} else {
		return projectUrl, validate(projectUrl)
	}
	return projectUrl, nil
}

// promptDestination prompts the user to select or provide a valid destination directory and validates the provided input.
// If no destination is provided, a selection prompt is displayed with preset destination options.
func promptDestination(dest string) (string, error) {
	validate := func(path string) error {
		invalidChars := `\:*?"<>|` // Characters not allowed in Unix paths

		// Check if the path starts with `/`
		if !strings.HasPrefix(path, "/") {
			return fmt.Errorf("invalid path: must start with '/'")
		}

		// Check for invalid characters
		for _, ch := range invalidChars {
			if strings.ContainsRune(path, ch) {
				return fmt.Errorf("invalid path: contains invalid character '%c'", ch)
			}
		}

		// Check if the path is empty
		if len(strings.TrimSpace(path)) == 0 {
			return fmt.Errorf("invalid path: cannot be empty")
		}

		// Path is valid
		return nil
	}

	if dest == "" {
		destinations := []string{
			"~/.local/share",
			"~/share",
			"/usr/local/share",
		}
		prompt := promptui.Select{
			Label: "Select Destination Directory",
			Items: destinations,
		}

		_, result, err := prompt.Run()
		if strings.HasPrefix("~", dest) {
			dest, err = utils.ExpandPath(result)
			if err != nil {
				return "", fmt.Errorf("Failed to expand path %v\n", err)
			}
		}
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}

	} else {
		return dest, validate(dest)
	}

	return dest, nil
}

// NewProject creates a new Project instance by prompting for missing or invalid inputs and initializing its fields.
func NewProject(args map[string]string) *Project {
	name := args["name"]
	projectUrl := args["url"]
	dest := args["dest"]
	command := args["command"]
	var err error
	name, err = promptName(name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	command, err = promptCommand(command)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	projectUrl, err = promptUrl(projectUrl)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	dest, err = promptDestination(dest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return &Project{
		Name:                 name,
		URL:                  projectUrl,
		DestinationDirectory: dest,
		DefaultCommand:       command,
	}
}
