package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Dalistor/gaver/pkg/types"
)

func Parse(initCommand *types.InitCommand) error {
	// Remover arquivos com extensão .tmplt_mysql, .tmplt_postgres e .tmplt_sqlite não utilizados
	switch initCommand.Database {
	case "postgres":
		removeFilesWithExtension(initCommand.Name, ".tmplt_mysql")
		removeFilesWithExtension(initCommand.Name, ".tmplt_sqlite")
		downloadSQLDriver("postgres")
	case "mysql":
		removeFilesWithExtension(initCommand.Name, ".tmplt_postgres")
		removeFilesWithExtension(initCommand.Name, ".tmplt_sqlite")
		downloadSQLDriver("mysql")
	case "sqlite":
		removeFilesWithExtension(initCommand.Name, ".tmplt_mysql")
		removeFilesWithExtension(initCommand.Name, ".tmplt_postgres")
		downloadSQLDriver("sqlite")
	}

	// Percorrer todas as pastas e arquivos do projeto e parsear arquivos com extensão .tmplt
	filepath.Walk(initCommand.Name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao percorrer pasta: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		if strings.Contains(info.Name(), ".") {
			extension := strings.Split(info.Name(), ".")[1]
			switch extension {
			case "tmplt":
				return parseFile(path, ".tmplt", ".go", map[string]string{
					"ProjectModuleName": initCommand.Name,
				})
			case "mod":
				return parseFile(path, ".mod", ".mod", map[string]string{
					"ProjectModuleName": initCommand.Name,
				})
			case "tmplt_mysql":
				return parseFile(path, ".tmplt_mysql", ".go", map[string]string{
					"ProjectModuleName": initCommand.Name,
				})
			case "tmplt_postgres":
				return parseFile(path, ".tmplt_postgres", ".go", map[string]string{
					"ProjectModuleName": initCommand.Name,
				})
			case "tmplt_sqlite":
				return parseFile(path, ".tmplt_sqlite", ".go", map[string]string{
					"ProjectModuleName": initCommand.Name,
				})
			default:
				return nil
			}
		}

		return nil
	})

	// Sincronizar arquivo mod
	fmt.Println("Sincronizando arquivo mod")
	if err := syncModFile(); err != nil {
		return fmt.Errorf("erro ao sincronizar arquivo mod: %w", err)
	}
	fmt.Println("Arquivo mod sincronizado com sucesso")

	return nil
}

func parseFile(filePath string, oldExt string, newExt string, vars map[string]string) error {
	fmt.Println("Parseando arquivo:", filePath)

	// Ler dados do arquivo
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	fileContentString := string(fileContent)
	// Substituir as variáveis no conteúdo do arquivo
	for key, value := range vars {
		fileContentString = strings.Replace(fileContentString, "{{."+key+"}}", value, -1)

		// Remover o arquivo template original
		os.Remove(filePath)

		// Salvar o arquivo com a nova extensão
		filePath = strings.Replace(filePath, oldExt, newExt, 1)
		err = os.WriteFile(filePath, []byte(fileContentString), 0644)
		if err != nil {
			return fmt.Errorf("erro ao salvar arquivo: %w", err)
		}
	}

	return nil
}

func removeFilesWithExtension(folder string, extension string) error {
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao percorrer pasta: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		if strings.Contains(info.Name(), extension) {
			return os.Remove(path)
		}

		return nil
	})

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

func syncModFile() error {
	err := exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		return fmt.Errorf("erro ao sincronizar arquivo mod: %w", err)
	}

	return nil
}
