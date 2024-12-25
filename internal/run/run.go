package run

import (
	"fmt"
	"gitlab.com/locke-codes/container-cli/internal/config"
	"gitlab.com/locke-codes/container-cli/internal/container"
	"log"
)

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
