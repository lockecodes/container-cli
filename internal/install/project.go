package install

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"gitlab.com/locke-codes/container-cli/internal/gitter"
	"gitlab.com/locke-codes/container-cli/internal/utils"
	"net/url"
	"os"
	"path"
	"strings"
)

type Project struct {
	Name                 string
	URL                  string
	DestinationDirectory string
}

func (p *Project) Path() string {
	return path.Join(p.DestinationDirectory, p.Name)
}

func (p *Project) Clone() error {
	client := gitter.NewGitter(p.Name, p.URL, p.Path())
	err := client.Clone()
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) Install() error {
	err := p.Uninstall()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	err = p.Clone()
	return err
}

func (p *Project) Uninstall() error {
	fmt.Printf("Removing %s\n", p.Path())
	err := os.RemoveAll(p.Path())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return nil
}

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

func NewProject(args map[string]string) *Project {
	name := args["name"]
	projectUrl := args["url"]
	dest := args["dest"]
	var err error
	name, err = promptName(name)
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
	}
}
