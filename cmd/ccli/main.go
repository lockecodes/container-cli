package main

import (
	"fmt"
	"github.com/locke-codes/container_cli/internal/update"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// version will be set during build
var version string

// main is the entry point of the application. It sets up the CLI interface with app configuration, commands, and flags.
func main() {
	app := &cli.App{
		Name:  "Container CLI",
		Usage: "Execute applications in containers",
		Commands: []*cli.Command{
			{
				Name:  "update",
				Usage: "Update to the latest version of the CLI",
				Action: func(c *cli.Context) error {
					update.Update()
					return nil
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

func GetVersion() string {
	return version
}
