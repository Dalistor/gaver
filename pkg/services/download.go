package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func DownloadTemplate(mode string, name string) error {
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

func downloadSQLDriver(sqlType string) error {
	switch sqlType {
	case "mysql":
		fmt.Println("Baixando driver de MySQL")
		exec.Command("go", "get", "gorm.io/driver/mysql").Run()
	case "postgres":
		fmt.Println("Baixando driver de PostgreSQL")
		exec.Command("go", "get", "gorm.io/driver/postgres").Run()
	case "sqlite":
		fmt.Println("Baixando driver de SQLite")
		exec.Command("go", "get", "gorm.io/driver/sqlite").Run()
	default:
		return fmt.Errorf("driver de SQL inválido")
	}

	return nil
}

func downloadFromGit(repoUrl string, branch string, path string) error {
	// criar pasta do projeto
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("erro ao criar pasta do projeto: %w", err)
	}

	// clonar repositório do template na branch específica
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           repoUrl,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
		Progress:      os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("erro ao clonar repositório da branch '%s' em '%s': %w", branch, path, err)
	}

	return nil
}

func DownloadSingleFileFromGit(repoFileUrl string, path string, fileName string) error {
	rawUrl := strings.Replace(repoFileUrl, "github.com", "raw.githubusercontent.com", 1)
	rawUrl = strings.Replace(rawUrl, "/blob/", "/", 1)

	// Criar diretório de destino se não existir
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório '%s': %w", path, err)
	}

	// Fazer requisição HTTP para baixar o arquivo
	resp, err := http.Get(rawUrl)
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição HTTP para '%s': %w", rawUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar arquivo: status code %d", resp.StatusCode)
	}

	// Criar arquivo de destino
	filePath := filepath.Join(path, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo '%s': %w", filePath, err)
	}
	defer out.Close()

	// Copiar conteúdo do response para o arquivo
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao escrever arquivo '%s': %w", filePath, err)
	}

	return nil
}
