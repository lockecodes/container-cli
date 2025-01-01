package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"gitlab.com/locke-codes/container-cli/internal/globals"
	"gitlab.com/locke-codes/container-cli/internal/install"
	"gitlab.com/locke-codes/container-cli/internal/run"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// version will be set during build
var version string

// main is the entry point of the application. It sets up the CLI interface with app configuration, commands, and flags.
func main() {
	cmd := &cli.Command{
		Name:  "Container CLI",
		Usage: "Execute applications in containers",
		Commands: []*cli.Command{
			{
				Name:  "install",
				Usage: "ContainerCLI the ccli binary",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return install.ContainerCLI()
				},
			},
			{
				Name:  "update",
				Usage: "Update to the latest version of the CLI",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					err := install.ContainerCLI()
					if err != nil {
						return err
					}
					ccliPath := filepath.Join(globals.HomeDir, ".local", "bin", "ccli")
					// Execute the command
					command := exec.Command(ccliPath, "version")
					command.Stdout = os.Stdout
					command.Stderr = os.Stderr

					// Run the command
					err = command.Run()
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Get the version of the CLI",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("Version: %s\n", version)
					return nil
				},
			},
			{
				Name:      "project",
				Usage:     "Ccli commands for projects",
				UsageText: "ccli project <command>",
				Commands: []*cli.Command{
					{
						Name:      "install",
						Usage:     "Install or update a new project",
						UsageText: "ccli project install",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "name",
								Usage: "name of the project",
							},
							&cli.StringFlag{
								Name:  "url",
								Usage: "Git url for the repository. E.g. ssh://git@gitlab.com/locke-codes/container-cli.git",
							},
							&cli.StringFlag{
								Name:  "dest",
								Usage: "destination directory for the project. E.g. ~/.local/share",
							},
							&cli.StringFlag{
								Name:  "command",
								Usage: "Default command to execute inside the container. E.g. 'bash'",
							},
							&cli.StringFlag{
								Name:  "alias",
								Usage: "Local command alias. E.g. 'bs' for 'big-salad' or 'hello' for 'hello-world",
							},
						},
						Action: func(ctx context.Context, cmd *cli.Command) error {
							args := map[string]string{
								"name":    cmd.String("name"),
								"url":     cmd.String("url"),
								"dest":    cmd.String("dest"),
								"command": cmd.String("command"),
								"alias":   cmd.String("alias"),
							}
							project := install.NewProject(args)
							fmt.Printf("Installing project: %s\nFrom url: %s\nTo directory: %s\n", project.Name, project.URL, project.DestinationDirectory)
							err := project.Install()
							if err != nil {
								panic(err)
							}
							return nil
						},
					},
					{
						Name:      "run",
						Usage:     "Run the project container",
						UsageText: "ccli project run",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							projectName := cmd.Args().First()
							args := cmd.Args().Tail()
							run.Run(projectName, args)
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
