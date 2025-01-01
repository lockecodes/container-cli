package install

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/utils"
	"gitlab.com/locke-codes/go-binary-updater/pkg/fileUtils"
	"gitlab.com/locke-codes/go-binary-updater/pkg/release"
)

// releaseObj represents an instance of the release.Release interface used for managing release lifecycle operations.
var releaseObj release.Release
var err error

func ExecContainerVersion() {
	ccliPath := filepath.Join(globals.HomeDir, ".local", "bin", "ccli")
	// Execute the command
	command := exec.Command(ccliPath, "version")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command
	_ = command.Run()
}

func ContainerCLIUpdate() error {
	engine, err := config.GetContainerEngineFromConfig()
	if err != nil {
		if err.Error() == "config file not found" {
			fmt.Printf(
				"No container engine found in config file. Running install instead.\n")
			return ContainerCLIInstall("", false)
		}
		return err
	}
	return ContainerCLI(engine)
}

func ContainerCLIInstall(engine string, force bool) error {
	if utils.FileExists(globals.DefaultContainerCliConfigPath) && !force {
		fmt.Printf("Config file already exists. Executing update instead.\n")
		return ContainerCLIUpdate()
	} else if utils.FileExists(globals.DefaultContainerCliConfigPath) && force {
		fmt.Printf("Config file already exists. Overwriting.\n")
		_ = os.Remove(globals.DefaultContainerCliConfigPath)
	}
	engine, err = promptEngine(engine)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Installing container-cli for %s\n", engine)
	return ContainerCLI(engine)
}
func ContainerCLI(engine string) error {
	// Define paths
	baseDir := filepath.Join(globals.HomeDir, ".local", "bin")
	fileConfig := fileUtils.FileConfig{
		VersionedDirectoryName: "container-cli",
		SourceBinaryName:       "container-cli",
		BinaryName:             "ccli",
		CreateGlobalSymlink:    false, // This isn't ready for use yet
		BaseBinaryDirectory:    baseDir,
		SourceArchivePath:      path.Join("/tmp", "container-cli.tar.gz"),
	}

	// Use GitLab implementation
	releaseObj = release.NewGitlabRelease(
		globals.ProjectId, // Ensure projectId matches the expected type
		fileConfig,
	)
	err = releaseObj.GetLatestRelease()
	if err != nil {
		return err
	}
	err = releaseObj.DownloadLatestRelease()
	if err != nil {
		return err
	}
	err = releaseObj.InstallLatestRelease()
	if err != nil {
		return err
	}

	configFile := config.NewContainerCliConfig(engine)
	err = configFile.LoadConfig()
	if !utils.FileExists(globals.DefaultContainerCliConfigPath) {
		err = configFile.SaveConfig()
		if err != nil {
			return err
		}
	}
	ExecContainerVersion()
	return nil
}

func promptEngine(engine string) (string, error) {
	validate := func(engineName string) error {
		if engineName != "" && engineName != "docker" && engineName != "podman" {
			return fmt.Errorf("invalid container engine: must be 'docker' or 'podman'")
		}
		// Path is valid
		return nil
	}

	if engine == "" {
		engines := []string{
			"docker",
			"podman",
		}
		prompt := promptui.Select{
			Label: "Select a container engine",
			Items: engines,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
		engine = result
	} else {
		fmt.Printf("Using container engine: %s\n", engine)
		return engine, validate(engine)
	}

	fmt.Printf("Using container engine: %s\n", engine)
	return engine, nil
}
