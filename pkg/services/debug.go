package services

import (
	"fmt"
	"os/exec"

	"github.com/Dalistor/gaver/pkg/types"
)

func Debug(initCommand *types.InitCommand) error {
	showGreetingMessage(initCommand)
	if initCommand.ProjectType == "api" {
		exec.Command("go", "run", "cmd/api/main.go").Run()
	}

	return nil
}

func showGreetingMessage(initCommand *types.InitCommand) {
	fmt.Println("===========================================")
	fmt.Println("Gaver - Framework multiplataformas")
	fmt.Println("Nome do projeto:", initCommand.Name)
	fmt.Println("Banco de dados:", initCommand.Database)
	fmt.Println("Tipo de projeto:", initCommand.ProjectType)
	fmt.Println("===========================================")
}
