package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	// Copiar templates de GitHub Actions se não existirem no repositório
	if err := copyGitHubActionsTemplates(repoPath); err != nil {
		return fmt.Errorf("erro ao copiar templates do GitHub Actions: %w", err)
	}

	return nil
}

func copyGitHubActionsTemplates(projectPath string) error {
	githubActionsDir := filepath.Join(projectPath, ".github", "workflows")

	// Verificar se o arquivo já existe no repositório
	existingFile := filepath.Join(githubActionsDir, "build.yml")
	if _, err := os.Stat(existingFile); err == nil {
		// Arquivo já existe, não precisa copiar
		return nil
	}

	// Verificar se o arquivo template já existe no repositório
	existingTemplate := filepath.Join(githubActionsDir, "build.yml.tmplt_yml")
	if _, err := os.Stat(existingTemplate); err == nil {
		// Template já existe no repositório, não precisa copiar
		return nil
	}

	// Criar diretório se não existir
	if err := os.MkdirAll(githubActionsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório .github/workflows: %w", err)
	}

	// Tentar encontrar o template local em diferentes locais possíveis
	possiblePaths := []string{
		"pkg/templates/github-actions/build.yml.tmplt_yml",
		filepath.Join("pkg", "templates", "github-actions", "build.yml.tmplt_yml"),
	}

	var templatePath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			templatePath = path
			break
		}
	}

	// Se não encontrou o template local, não há problema (será incluído no repositório remoto)
	if templatePath == "" {
		return nil
	}

	// Copiar template para o projeto
	destPath := filepath.Join(githubActionsDir, "build.yml.tmplt_yml")
	if err := copyFile(templatePath, destPath); err != nil {
		return fmt.Errorf("erro ao copiar template: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
