package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/utils"
	"log"
	"path"
)

type ProjectConfig struct {
	Name           string `koanf:"name"`
	Path           string `koanf:"path"`
	Dockerfile     string `koanf:"dockerfile"`
	BuildDirectory string `koanf:"buildDirectory"`
	BuildContext   string `koanf:"buildContext"`
	DefaultCommand string `koanf:"defaultCommand"`
}

type ContainerCliConfig struct {
	ContainerEngine string          `koanf:"containerEngine"`
	Path            string          `koanf:"path"`
	Projects        []ProjectConfig `koanf:"projects"`
}

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var (
	k      = koanf.New(".")
	parser = yaml.Parser()
	err    error
)

func NewContainerCliConfig() *ContainerCliConfig {
	configFile := ContainerCliConfig{
		ContainerEngine: "podman",
		Path:            globals.DefaultContainerCliConfigPath,
		Projects:        []ProjectConfig{},
	}
	configFile.KoanfLoad()
	return &configFile
}

func (c *ContainerCliConfig) KoanfLoad() {
	_ = k.Load(structs.Provider(c, "koanf"), nil)
}

// LoadConfig reads a YAML file specified by filename and unmarshals its content into a ContainerCliConfig struct.
// It returns the loaded ContainerCliConfig and any error encountered during the file reading or unmarshalling process.
func (c *ContainerCliConfig) LoadConfig() error {
	if !utils.FileExists(c.Path) {
		return nil
	}
	log.Printf("Loading config from %s", c.Path)
	log.Printf("Parsing YAML")
	if err = k.Load(file.Provider(c.Path), parser); err != nil {
		return fmt.Errorf("error reading %s: %w", c.Path, err)
	}
	var config ContainerCliConfig
	log.Printf("Unmarshalling YAML")
	if err = k.Unmarshal("", &config); err != nil {
		return fmt.Errorf("error parsing YAML: %w", err)
	}
	c.ContainerEngine = config.ContainerEngine
	c.Projects = config.Projects
	return nil
}

func (c *ContainerCliConfig) SaveConfig() error {
	// Marshal the instance back to JSON.
	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	c.KoanfLoad()
	marshalledBytes, err := k.Marshal(parser)
	if err != nil {
		return err
	}
	// Write the byte slice to the file
	if err = utils.WriteToFile(c.Path, marshalledBytes); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	} else {
		fmt.Printf("Data successfully written to %s\n", c.Path)
	}
	return nil
}

func (c *ContainerCliConfig) GetProject(name string) *ProjectConfig {
	for _, project := range c.Projects {
		if project.Name == name {
			return &project
		}
	}
	return nil
}

// Helper function to copy a slice
func copySlice(original []ProjectConfig) []ProjectConfig {
	// Create a new slice with the same length as the original
	copied := make([]ProjectConfig, len(original))
	copy(copied, original) // Copies the data from the original slice to the new slice
	return copied
}

// ReplaceProjectByName replaces a person in the slice by their name
func (c *ContainerCliConfig) ReplaceProjectByName(name string, newProject ProjectConfig) error {
	fmt.Printf("Replacing project %s with %s\n", name, newProject.Name)
	projectList := copySlice(c.Projects)
	for i, project := range projectList {
		if project.Name == name {
			c.Projects[i] = newProject
			return nil
		}
	}
	return fmt.Errorf("project with name %s not found", name)
}

func (c *ContainerCliConfig) GetProjectPath(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return project.Path
	}
	return ""
}

func (c *ContainerCliConfig) GetProjectBuildContext(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return project.BuildContext
	}
	return ""
}

func (c *ContainerCliConfig) GetProjectBuildDirectory(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return project.BuildDirectory
	}
	return ""
}

func (c *ContainerCliConfig) GetProjectDockerfile(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return project.Dockerfile
	}
	return ""
}

func (c *ContainerCliConfig) GetProjectDockerfilePath(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return path.Join(project.Path, project.Dockerfile)
	}
	return ""
}
