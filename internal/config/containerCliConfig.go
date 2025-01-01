package config

import (
	"fmt"
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/utils"
)

// ContainerCliConfig represents the configuration for the container CLI, including engine, path, and project details.
type ContainerCliConfig struct {
	ContainerEngine string          `koanf:"containerEngine"`
	Path            string          `koanf:"path"`
	Projects        []ProjectConfig `koanf:"projects"`
}

var (
	// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	k      = koanf.New(".")
	parser = yaml.Parser()
	err    error
)

// NewContainerCliConfig initializes a default configuration for the container CLI and loads the default settings
// directly into koanf
func NewContainerCliConfig(engine string) *ContainerCliConfig {
	configFile := ContainerCliConfig{
		ContainerEngine: engine,
		Path:            globals.DefaultContainerCliConfigPath,
		Projects:        []ProjectConfig{},
	}
	configFile.KoanfLoad()
	return &configFile
}

// LoadConfig loads a container CLI configuration file into a ContainerCliConfig instance and returns it or an error.
func LoadConfig() (*ContainerCliConfig, error) {
	if !utils.FileExists(globals.DefaultContainerCliConfigPath) {
		return nil, fmt.Errorf("config file not found")
	}
	configFile := ContainerCliConfig{
		Path: globals.DefaultContainerCliConfigPath,
	}
	err := configFile.LoadConfig()
	if err != nil {
		return nil, err
	}
	return &configFile, nil
}

// GetContainerEngineFromConfig retrieves the container engine from the configuration file or returns an error if not set.
func GetContainerEngineFromConfig() (string, error) {
	configFile, err := LoadConfig()
	if err != nil {
		return "", err
	}
	fmt.Printf("Container engine: %s\n", configFile.ContainerEngine)
	if configFile.ContainerEngine == "" {
		return "", fmt.Errorf("containerEngine not set")
	}
	return configFile.ContainerEngine, nil
}

// KoanfLoad loads the ContainerCliConfig struct into the global koanf instance using the "koanf" struct tags.
func (c *ContainerCliConfig) KoanfLoad() {
	_ = k.Load(structs.Provider(c, "koanf"), nil)
}

// LoadConfig reads a YAML file specified by filename and unmarshals its content into a ContainerCliConfig struct.
// It returns the loaded ContainerCliConfig and any error encountered during the file reading or unmarshalling process.
func (c *ContainerCliConfig) LoadConfig() error {
	// If the config file doesn't exist, just return
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

// SaveConfig saves the ContainerCliConfig instance to a file in YAML format by marshalling it and writing to the
// specified path.
func (c *ContainerCliConfig) SaveConfig() error {
	// Marshal the instance back to YAML.
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

// GetProject retrieves the ProjectConfig for the given project name from the ContainerCliConfig.
// Returns nil if the project is not found.
func (c *ContainerCliConfig) GetProject(name string) *ProjectConfig {
	for _, project := range c.Projects {
		if project.Name == name {
			return &project
		}
	}
	return nil
}

// ReplaceProjectByName replaces a person in the slice by their name
func (c *ContainerCliConfig) ReplaceProjectByName(name string, newProject ProjectConfig) error {
	fmt.Printf("Replacing project %s with %s\n", name, newProject.Name)
	projectList := utils.CopySlice(c.Projects)
	for i, project := range projectList {
		if project.Name == name {
			c.Projects[i] = newProject
			return nil
		}
	}
	return fmt.Errorf("project with name %s not found", name)
}

// GetProjectPath returns the file system path of the specified project name if found, otherwise it returns an empty string.
func (c *ContainerCliConfig) GetProjectPath(name string) string {
	project := c.GetProject(name)
	if project != nil {
		return project.Path
	}
	return ""
}
