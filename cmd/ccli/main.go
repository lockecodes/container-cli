package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gitlab.com/locke-codes/go-binary-updater/pkg/release"
	"log"
	"os"
)

// version will be set during build
var version string
var releaseObj release.Release

const projectId int = 47137983

func install() error {
	fileConfig := release.FileConfig{
		VersionedDirectoryName: "container-cli",
		SourceBinaryName:       "container-cli",
		BinaryName:             "ccli",
		CreateGlobalSymlink:    false, // This isn't ready for use yet
	}

	// Use GitLab implementation
	releaseObj = release.NewGitlabRelease(
		projectId, // Ensure projectId matches the expected type
		fileConfig,
	)
	err := releaseObj.GetLatestRelease()
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

// main is the entry point of the application. It sets up the CLI interface with app configuration, commands, and flags.
func main() {
	app := &cli.App{
		Name:  "Container CLI",
		Usage: "Execute applications in containers",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "Install the ccli binary",
				Action: func(c *cli.Context) error {
					return install()
				},
			},
			{
				Name:  "update",
				Usage: "Update to the latest version of the CLI",
				Action: func(c *cli.Context) error {
					return install()
				},
			},
			{
				Name:  "version",
				Usage: "Get the version of the CLI",
				Action: func(c *cli.Context) error {
					fmt.Printf("Version: %s\n", version)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
