package install

import (
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/utils"
	"gitlab.com/locke-codes/go-binary-updater/pkg/fileUtils"
	"gitlab.com/locke-codes/go-binary-updater/pkg/release"
	"path"
	"path/filepath"
)

// releaseObj represents an instance of the release.Release interface used for managing release lifecycle operations.
var releaseObj release.Release
var err error

// ContainerCLI handles the installation and configuration process for the Container CLI, including release management.
func ContainerCLI() error {
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

	configFile := config.NewContainerCliConfig()
	err = configFile.LoadConfig()
	if !utils.FileExists(globals.DefaultContainerCliConfigPath) {
		err = configFile.SaveConfig()
		if err != nil {
			return err
		}
	}
	return nil
}
