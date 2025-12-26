package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Dalistor/gaver/pkg/types"
)

func ParseInit(initCommand *types.InitCommand) error {
	// Remover arquivos com extensão .tmplt_mysql, .tmplt_postgres e .tmplt_sqlite não utilizados
	switch initCommand.Database {
	case "postgres":
		removeFilesWithExtension(initCommand.Name, ".tmplt_mysql")
		removeFilesWithExtension(initCommand.Name, ".tmplt_sqlite")
		removeFilesWithExtension(initCommand.Name, ".db")
		downloadSQLDriver("postgres")
	case "mysql":
		removeFilesWithExtension(initCommand.Name, ".tmplt_postgres")
		removeFilesWithExtension(initCommand.Name, ".tmplt_sqlite")
		removeFilesWithExtension(initCommand.Name, ".db")
		downloadSQLDriver("mysql")
	case "sqlite":
		removeFilesWithExtension(initCommand.Name, ".tmplt_mysql")
		removeFilesWithExtension(initCommand.Name, ".tmplt_postgres")
		downloadSQLDriver("sqlite")
	}

	switch initCommand.ProjectType {
	case "api", "web", "desktop":
		removeFilesWithExtension(initCommand.Name, ".tmplt_desktop_yml")
		removeFilesWithExtension(initCommand.Name, ".tmplt_mobile_yml")
	case "mobile":
		removeFilesWithExtension(initCommand.Name, ".tmplt_desktop_yml")
		removeFilesWithExtension(initCommand.Name, ".tmplt_api_yml")
	}

	// Percorrer todas as pastas e arquivos do projeto e parsear arquivos com extensão .tmplt
	filepath.Walk(initCommand.Name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao percorrer pasta: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		fileName := info.Name()
		vars := map[string]string{
			"ProjectModuleName": initCommand.Name,
		}

		// Verificar extensões compostas primeiro (ordem importa)
		switch {
		case strings.HasSuffix(fileName, ".tmplt_desktop_yml"):
			return parseFile(path, ".tmplt_desktop_yml", ".yml", vars)
		case strings.HasSuffix(fileName, ".tmplt_mobile_yml"):
			return parseFile(path, ".tmplt_mobile_yml", ".yml", vars)
		case strings.HasSuffix(fileName, ".tmplt_api_yml"):
			return parseFile(path, ".tmplt_api_yml", ".yml", vars)
		case strings.HasSuffix(fileName, ".tmplt_yml"):
			return parseFile(path, ".tmplt_yml", ".yml", vars)
		case strings.HasSuffix(fileName, ".tmplt_yaml"):
			return parseFile(path, ".tmplt_yaml", ".yaml", vars)
		case strings.HasSuffix(fileName, ".tmplt_mysql"):
			return parseFile(path, ".tmplt_mysql", ".go", vars)
		case strings.HasSuffix(fileName, ".tmplt_postgres"):
			return parseFile(path, ".tmplt_postgres", ".go", vars)
		case strings.HasSuffix(fileName, ".tmplt_sqlite"):
			return parseFile(path, ".tmplt_sqlite", ".go", vars)
		case strings.HasSuffix(fileName, ".tmplt"):
			return parseFile(path, ".tmplt", ".go", vars)
		case strings.HasSuffix(fileName, ".mod"):
			return parseFile(path, ".mod", ".mod", vars)
		default:
			return nil
		}
	})

	// Sincronizar arquivo mod
	fmt.Println("Sincronizando arquivo mod")
	if err := syncModFile(initCommand.Database, initCommand.Name); err != nil {
		return fmt.Errorf("erro ao sincronizar arquivo mod: %w", err)
	}
	fmt.Println("Arquivo mod sincronizado com sucesso")

	// Setar arquivo module
	fmt.Println("Setando arquivo module")
	if err := setGaverModuleFile(initCommand); err != nil {
		return fmt.Errorf("erro ao setar arquivo module: %w", err)
	}
	fmt.Println("Arquivo module setado com sucesso")

	// Remover pasta .git
	if err := removeDotGitFolder(initCommand.Name); err != nil {
		return fmt.Errorf("erro ao remover pasta .git: %w", err)
	}

	return nil
}

func ParseModule(moduleCommand *types.ModuleCommand, projectName string) error {
	modulePath := fmt.Sprintf("modules/%s", moduleCommand.Name)

	// Andar pelos arquivos do módulo e personalizar eles
	filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao percorrer pasta: %w", err)
		}
		if info.IsDir() {
			return nil
		}

		// Remover arquivos que não são necessários para o tipo de módulo e converter para .go
		var vars = map[string]string{
			"ProjectModuleName": projectName,
			"ModuleName":        moduleCommand.Name,
			"ModelName":         moduleCommand.Controller,
		}

		fileType := strings.Split(info.Name(), ".")[0]

		if vars["ModelName"] == "" {
			// Se não houver model name, então é um download genérico de templates
			switch moduleCommand.Type {
			case "crud":
				removeFilesWithExtension(path, ".tmplt_service")

				switch fileType {
				case "handler":
					removeFile(path)
				case "repository":
					removeFile(path)
				case "service":
					removeFile(path)
				}

				switch {
				case strings.HasSuffix(info.Name(), ".tmplt_crud"):
					parseFile(path, ".tmplt_crud", ".go", vars)
				case strings.HasSuffix(info.Name(), ".tmplt"):
					parseFile(path, ".tmplt", ".go", vars)
				}

			case "service":
				removeFilesWithExtension(path, ".tmplt_crud")

				switch {
				case strings.HasSuffix(info.Name(), ".tmplt_service"):
					parseFile(path, ".tmplt_service", ".go", vars)
				case strings.HasSuffix(info.Name(), ".tmplt"):
					parseFile(path, ".tmplt", ".go", vars)
				}
			}
		} else {
			// Se houver model name, então é um download específico de templates CRUD
			parseFile(path, ".tmplt_crud", ".go", vars)

			// Renomear arquivo models para o nome do model
			if info.Name() == "models.go" {
				os.Rename(path, fmt.Sprintf("%s.go", strings.ToLower(vars["ModelName"])))
			}
		}

		return nil
	})

	return nil
}

func parseFile(filePath string, oldExt string, newExt string, vars map[string]string) error {
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

func removeFile(path string) error {
	return os.Remove(path)
}

func syncModFile(database string, folder string) error {
	os.Chdir(folder)

	if err := exec.Command("go", "get", fmt.Sprintf("gorm.io/driver/%s", database)).Run(); err != nil {
		return fmt.Errorf("erro ao sincronizar arquivo mod: %w", err)
	}

	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return fmt.Errorf("erro ao sincronizar arquivo mod: %w", err)
	}

	os.Chdir("..")

	return nil
}

func setGaverModuleFile(initCommand *types.InitCommand) error {
	moduleFile := &types.GaverModuleFile{
		Type:                initCommand.ProjectType,
		ProjectName:         initCommand.Name,
		ProjectVersion:      "1.0.0",
		ProjectModules:      []string{},
		ProjectDatabaseType: initCommand.Database,
		MigrationTag:        0,
	}

	jsonData, err := json.Marshal(moduleFile)
	if err != nil {
		return fmt.Errorf("erro ao serializar arquivo module: %w", err)
	}

	os.WriteFile(fmt.Sprintf("%s/gaverModule.json", initCommand.Name), jsonData, 0644)

	return nil
}

func removeDotGitFolder(folder string) error {
	os.RemoveAll(fmt.Sprintf("%s/.git", folder))
	return nil
}
