package services

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Dalistor/gaver/pkg/types"
)

func MakeMigrations(gaverModuleFile *types.GaverModuleFile) error {
	cmd := exec.Command("go", "run", "cmd/api/main.go", "makemigrations")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar as migrações, certifique-se de que você está na pasta do projeto e que o caminho cmd/api/main.go existe")
	}

	return nil
}

func Migrate(gaverModuleFile *types.GaverModuleFile) error {
	cmd := exec.Command("go", "run", "cmd/api/main.go", "migrate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao executar as migrações, certifique-se de que você está na pasta do projeto e que o caminho cmd/api/main.go existe")
	}

	return nil
}
