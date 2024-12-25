package run

import (
	"fmt"
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/container"
	"log"
)

// Run initializes and runs a container for the specified project with the provided arguments.
// It loads the project configuration, builds the container, and runs it with specified parameters.
// Logs any errors encountered during configuration loading, building, or running of the container.
// TODO: This method of running the containers is still incomplete
//
//	Because the cli library doesn't support passing arbitrary flags it is unable to fully support calling the
//	commands. This has been mitigated by modifying the script `~/.local/bin/command` to include the fully bash
//	command for running the container. This actually may be sufficient which would mean it would make sense to
//	just remove this functionality. Time will tell
func Run(projectName string, args []string) {
	fmt.Printf("Running %s with args %v\n", projectName, args)
	configFile := config.NewContainerCliConfig()
	err := configFile.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	project := configFile.GetProject(projectName)
	if project == nil {
		log.Fatal("Project not found")
	}
	containerObj := container.NewPodman(project)
	err = containerObj.Build()
	if err != nil {
		log.Fatal(err)
	}
	err = containerObj.Run(args)
	if err != nil {
		log.Fatal(err)
	}
}
