package services

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func DownloadAPI(mode string, name string) error {
	switch mode {
	case "api":
		return downloadAPITemplate(name)
	case "module":
		return downloadModuleTemplate(name)
	}

	return fmt.Errorf("modo inválido")
}

func downloadAPITemplate(name string) error {
	if err := downloadFromGit("https://github.com/Dalistor/Gaver-Modules", "api", name); err != nil {
		return fmt.Errorf("erro ao baixar template da API: %w", err)
	}

	return nil
}

func downloadModuleTemplate(name string) error {
	path := fmt.Sprintf("modules/%s", name)

	if err := downloadFromGit("https://github.com/Dalistor/Gaver-Modules", "module", path); err != nil {
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
