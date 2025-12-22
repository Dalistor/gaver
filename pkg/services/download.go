package services

import (
	"fmt"
	"os"

	"github.com/Dalistor/gaver/pkg/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Download(initCommand *types.InitCommand) error {
	switch initCommand.ProjectType {
	case "api":
		return downloadAPITemplate(initCommand)
	}

	return nil
}

func downloadAPITemplate(initCommand *types.InitCommand) error {
	// criar pasta do projeto
	repoPath := initCommand.Name
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar pasta do projeto: %w", err)
	}

	// clonar repositório do template na branch específica
	repoUrl := "https://github.com/Dalistor/Gaver-Modules"
	branch := "api"

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
