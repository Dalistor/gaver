package structure

import (
	"fmt"
	"os"
	"path/filepath"

	templates "github.com/Dalistor/gaver/internal/templates"
)

// CreateProjectFolders cria a estrutura de pastas do projeto
func CreateProjectFolders(projectName string) error {
	dirs := []string{
		projectName,
		filepath.Join(projectName, "cmd", "server"),
		filepath.Join(projectName, "config", "env"),
		filepath.Join(projectName, "config", "middlewares"),
		filepath.Join(projectName, "config", "cors"),
		filepath.Join(projectName, "config", "database"),
		filepath.Join(projectName, "config", "database", "migrations"),
		filepath.Join(projectName, "config", "routines"),
		filepath.Join(projectName, "config", "routes"),
		filepath.Join(projectName, "config", "modules"),
		filepath.Join(projectName, "modules"),
		filepath.Join(projectName, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

type ProjectConfig struct {
	ProjectName          string
	DatabaseDriver       string
	DatabaseDriverImport string
	DatabasePort         string
	DatabaseUser         string
}

// GenerateInitialFiles gera arquivos iniciais do projeto
func GenerateInitialFiles(projectName, database string) error {
	// Criar estrutura de pastas primeiro
	if err := CreateProjectFolders(projectName); err != nil {
		return fmt.Errorf("erro ao criar pastas: %w", err)
	}

	config := ProjectConfig{
		ProjectName:          projectName,
		DatabaseDriver:       getDatabaseDriver(database),
		DatabaseDriverImport: getDatabaseDriverImport(database),
		DatabasePort:         getDatabasePort(database),
		DatabaseUser:         getDatabaseUser(database),
	}

	gen := templates.New(projectName)

	// Gerar arquivos de config
	files := map[string]string{
		"config_env.tmpl":         "config/env/env.go",
		"config_middlewares.tmpl": "config/middlewares/middlewares.go",
		"config_cors.tmpl":        "config/cors/cors.go",
		"config_database.tmpl":    "config/database/database.go",
		"migration_table.tmpl":    "config/database/migrations/migrations.go",
		"routines.tmpl":           "config/routines/routines.go",
		"config_routes.tmpl":      "config/routes/routes.go",
		"config_modules.tmpl":     "config/modules/modules.go",
		"main.tmpl":               "cmd/server/main.go",
		"env.tmpl":                ".env",
		"gitignore.tmpl":          ".gitignore",
		"go_mod.tmpl":             "go.mod",
		"readme.tmpl":             "README.md",
	}

	for template, output := range files {
		if err := gen.Generate(template, output, config); err != nil {
			return fmt.Errorf("erro ao gerar %s: %w", output, err)
		}
	}

	return nil
}

func getDatabaseDriver(db string) string {
	drivers := map[string]string{
		"postgres": "postgres",
		"mysql":    "mysql",
		"sqlite":   "sqlite",
	}
	if driver, ok := drivers[db]; ok {
		return driver
	}
	return "mysql"
}

func getDatabaseDriverImport(db string) string {
	driver := getDatabaseDriver(db)
	imports := map[string]string{
		"postgres": "gorm.io/driver/postgres",
		"mysql":    "gorm.io/driver/mysql",
		"sqlite":   "gorm.io/driver/sqlite",
	}
	return imports[driver]
}

func getDatabasePort(db string) string {
	ports := map[string]string{
		"postgres": "5432",
		"mysql":    "3306",
		"sqlite":   "",
	}
	return ports[db]
}

func getDatabaseUser(db string) string {
	users := map[string]string{
		"postgres": "postgres",
		"mysql":    "root",
		"sqlite":   "",
	}
	return users[db]
}
