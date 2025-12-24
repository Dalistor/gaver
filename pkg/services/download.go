package services

import (
	"fmt"
	"os"

	"github.com/Dalistor/gaver/pkg/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Download(initCommand *types.InitCommand, moduleCommand *types.GaverModuleFile) error {
	if initCommand != nil {
		switch initCommand.ProjectType {
		case "api":
			return downloadAPITemplate(initCommand)
		}
	}

	if moduleCommand != nil {
		downloadModuleTemplate(moduleCommand)
	}

	return nil
}

func downloadAPITemplate(initCommand *types.InitCommand) error {
	if err := downloadFromGit("https://github.com/Dalistor/Gaver-Modules", "api", initCommand.Name); err != nil {
		return fmt.Errorf("erro ao baixar template da API: %w", err)
	}

	return nil
}

func downloadModuleTemplate(moduleCommand *types.GaverModuleFile) error {
	if err := downloadFromGit("https://github.com/Dalistor/Gaver-Modules", "module", moduleCommand.ProjectName); err != nil {
		return fmt.Errorf("erro ao baixar template do módulo: %w", err)
	}

	return nil
}

func downloadFromGit(repoUrl string, branch string, repoPath string) error {
	// criar pasta do projeto
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar pasta do projeto: %w", err)
	}

	// clonar repositório do template na branch específica
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("erro ao clonar repositório da branch '%s' em '%s': %w", branch, repoPath, err)
	}

	return nil
}
