package install

import (
	"fmt"
	"gitlab.com/locke-codes/go-binary-updater/pkg/fileUtils"
	"gitlab.com/locke-codes/go-binary-updater/pkg/release"
	"os"
	"path"
	"path/filepath"
)

var releaseObj release.Release

const projectId string = "47137983"

func ContainerCLI() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Define paths
	baseDir := filepath.Join(homeDir, ".local", "bin")
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
		projectId, // Ensure projectId matches the expected type
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
	return nil
}
