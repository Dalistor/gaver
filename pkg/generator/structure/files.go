package structure

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	templates "github.com/Dalistor/gaver/internal/templates"
	"github.com/Dalistor/gaver/pkg/config"
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
	ProjectType          string
	ServerPort           string
}

// GenerateInitialFiles gera arquivos iniciais do projeto
func GenerateInitialFiles(projectName, database, projectType string) error {
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
		ProjectType:          projectType,
		ServerPort:           "8080", // Porta padrão do servidor
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
		"env_example.tmpl":        ".env.example",
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
		"sqlite":   "github.com/glebarez/sqlite", // Driver SQLite puro Go, não requer CGO
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

// FrontendConfig representa a configuração para gerar frontend
type FrontendConfig struct {
	ProjectName string
	ServerPort  string
}

// GenerateAndroidFrontend gera a estrutura frontend para Android
func GenerateAndroidFrontend(projectName string, projectConfig *config.ProjectConfig) error {
	gen := templates.New(projectName)

	frontendConfig := FrontendConfig{
		ProjectName: projectName,
		ServerPort:  projectConfig.ServerPort,
	}

	// Criar estrutura de pastas frontend
	frontendDirs := []string{
		filepath.Join(projectName, "frontend", "src", "composables"),
		filepath.Join(projectName, "frontend", "src", "api"),
		filepath.Join(projectName, "frontend", "src", "components"),
		filepath.Join(projectName, "frontend", "src", "pages"),
		filepath.Join(projectName, "frontend", "src", "layouts"),
		filepath.Join(projectName, "frontend", "src", "router"),
		filepath.Join(projectName, "frontend", "src", "boot"),
		filepath.Join(projectName, "frontend", "src", "assets"),
		filepath.Join(projectName, "frontend", "src", "css"),
	}

	for _, dir := range frontendDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar arquivos frontend
	frontendFiles := map[string]string{
		"quasar_config.tmpl":            "frontend/quasar.config.js",
		"package_json_android.tmpl":     "frontend/package.json",
		"capacitor_config.tmpl":         "frontend/capacitor.config.js",
		"frontend_env.tmpl":             "frontend/.env",
		"composable_api.tmpl":           "frontend/src/composables/useApi.ts",
		"api_client.tmpl":               "frontend/src/api/client.js",
		"router_config.tmpl":            "frontend/src/router/index.js",
		"router_routes.tmpl":            "frontend/src/router/routes.js",
		"app_main.tmpl":                 "frontend/src/main.js",
		"app_vue.tmpl":                  "frontend/src/App.vue",
		"index_html.tmpl":               "frontend/index.html",
		"layout_main.tmpl":              "frontend/src/layouts/MainLayout.vue",
		"page_index.tmpl":               "frontend/src/pages/IndexPage.vue",
		"page_error.tmpl":               "frontend/src/pages/ErrorNotFound.vue",
		"component_essential_link.tmpl": "frontend/src/components/EssentialLink.vue",
		"boot_axios.tmpl":               "frontend/src/boot/axios.js",
		"app_scss.tmpl":                 "frontend/src/css/app.scss",
	}

	for template, output := range frontendFiles {
		if err := gen.Generate(template, output, frontendConfig); err != nil {
			return fmt.Errorf("erro ao gerar %s: %w", output, err)
		}
	}

	// Instalar dependências npm (inclui Capacitor)
	fmt.Println("📦 Instalando dependências npm (incluindo Capacitor)...")
	frontendPath := filepath.Join(projectName, "frontend")

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	// Mudar para diretório frontend
	if err := os.Chdir(frontendPath); err != nil {
		return fmt.Errorf("erro ao mudar para diretório frontend: %w", err)
	}

	npmCmd := exec.Command("npm", "install")
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr

	if err := npmCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao instalar dependências npm. Execute 'npm install' manualmente no diretório frontend.")
	} else {
		fmt.Println("✓ Dependências npm instaladas (Capacitor incluído)")
	}

	// Inicializar Capacitor após instalar dependências
	fmt.Println("🔧 Inicializando Capacitor...")
	capacitorInitCmd := exec.Command("npx", "cap", "init", projectName, "--web-dir=dist")
	capacitorInitCmd.Stdout = os.Stdout
	capacitorInitCmd.Stderr = os.Stderr
	if err := capacitorInitCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao inicializar Capacitor. Execute 'npx cap init' manualmente se necessário.")
	}

	// Adicionar plataforma Android
	fmt.Println("📱 Adicionando plataforma Android...")
	capacitorAddCmd := exec.Command("npx", "cap", "add", "android")
	capacitorAddCmd.Stdout = os.Stdout
	capacitorAddCmd.Stderr = os.Stderr
	if err := capacitorAddCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao adicionar Android. Execute 'npx cap add android' manualmente se necessário.")
	}

	// Voltar para diretório original
	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("erro ao voltar para diretório original: %w", err)
	}

	return nil
}

// GenerateDesktopFrontend gera a estrutura frontend para Desktop
func GenerateDesktopFrontend(projectName string, projectConfig *config.ProjectConfig) error {
	gen := templates.New(projectName)

	frontendConfig := FrontendConfig{
		ProjectName: projectName,
		ServerPort:  projectConfig.ServerPort,
	}

	// Criar estrutura de pastas frontend
	frontendDirs := []string{
		filepath.Join(projectName, "frontend", "src", "composables"),
		filepath.Join(projectName, "frontend", "src", "api"),
		filepath.Join(projectName, "frontend", "src", "components"),
		filepath.Join(projectName, "frontend", "src", "pages"),
		filepath.Join(projectName, "frontend", "src", "layouts"),
		filepath.Join(projectName, "frontend", "src", "router"),
		filepath.Join(projectName, "frontend", "src", "boot"),
		filepath.Join(projectName, "frontend", "src", "assets"),
		filepath.Join(projectName, "frontend", "src", "css"),
		filepath.Join(projectName, "frontend", "src-electron"),
	}

	for _, dir := range frontendDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar arquivos frontend
	frontendFiles := map[string]string{
		"quasar_config.tmpl":            "frontend/quasar.config.js",
		"package_json_desktop.tmpl":     "frontend/package.json",
		"frontend_env.tmpl":             "frontend/.env",
		"electron_main.tmpl":            "frontend/src-electron/electron-main.js",
		"composable_api.tmpl":           "frontend/src/composables/useApi.ts",
		"api_client.tmpl":               "frontend/src/api/client.js",
		"router_config.tmpl":            "frontend/src/router/index.js",
		"router_routes.tmpl":            "frontend/src/router/routes.js",
		"app_main.tmpl":                 "frontend/src/main.js",
		"app_vue.tmpl":                  "frontend/src/App.vue",
		"index_html.tmpl":               "frontend/index.html",
		"layout_main.tmpl":              "frontend/src/layouts/MainLayout.vue",
		"page_index.tmpl":               "frontend/src/pages/IndexPage.vue",
		"page_error.tmpl":               "frontend/src/pages/ErrorNotFound.vue",
		"component_essential_link.tmpl": "frontend/src/components/EssentialLink.vue",
		"boot_axios.tmpl":               "frontend/src/boot/axios.js",
		"app_scss.tmpl":                 "frontend/src/css/app.scss",
	}

	for template, output := range frontendFiles {
		if err := gen.Generate(template, output, frontendConfig); err != nil {
			return fmt.Errorf("erro ao gerar %s: %w", output, err)
		}
	}

	// Copiar logo.png para assets do projeto
	assetsPath := filepath.Join(projectName, "frontend", "src", "assets")
	logoDest := filepath.Join(assetsPath, "logo.png")
	
	// Obter diretório atual para tentar encontrar o logo
	currentDir, _ := os.Getwd()
	
	// Tentar encontrar o logo em diferentes locais
	possibleLogoPaths := []string{
		filepath.Join("assets", "logo.png"),                    // Diretório atual (desenvolvimento)
		filepath.Join("..", "assets", "logo.png"),              // Um nível acima
		filepath.Join("..", "..", "assets", "logo.png"),       // Dois níveis acima
		filepath.Join(currentDir, "assets", "logo.png"),        // Diretório atual absoluto
	}
	
	var logoSource string
	for _, possiblePath := range possibleLogoPaths {
		if _, err := os.Stat(possiblePath); err == nil {
			logoSource = possiblePath
			break
		}
	}
	
	// Se encontrou o logo, copiar
	if logoSource != "" {
		sourceFile, err := os.Open(logoSource)
		if err == nil {
			defer sourceFile.Close()
			
			destFile, err := os.Create(logoDest)
			if err == nil {
				defer destFile.Close()
				
				_, err = io.Copy(destFile, sourceFile)
				if err == nil {
					fmt.Println("✓ Logo copiado para assets/")
				} else {
					fmt.Printf("⚠️  Aviso: Erro ao copiar logo: %v\n", err)
				}
			}
		}
	} else {
		fmt.Println("ℹ️  Logo não encontrado - será necessário adicionar src/assets/logo.png manualmente")
		fmt.Println("   Você pode copiar o logo de assets/logo.png do framework ou usar seu próprio logo")
	}

	// Instalar dependências npm (inclui Electron)
	fmt.Println("📦 Instalando dependências npm (incluindo Electron)...")
	frontendPath := filepath.Join(projectName, "frontend")

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	// Mudar para diretório frontend
	if err := os.Chdir(frontendPath); err != nil {
		return fmt.Errorf("erro ao mudar para diretório frontend: %w", err)
	}

	npmCmd := exec.Command("npm", "install")
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr

	if err := npmCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao instalar dependências npm. Execute 'npm install' manualmente no diretório frontend.")
	} else {
		fmt.Println("✓ Dependências npm instaladas (Electron incluído)")
	}

	// Gerar ícones do Electron usando @quasar/icongenie (CLI standalone)
	// Verificar se o logo existe antes de tentar gerar ícones
	logoPath := filepath.Join("src", "assets", "logo.png")
	if _, err := os.Stat(logoPath); err == nil {
		fmt.Println("🎨 Gerando ícones do Electron...")
		// @quasar/icongenie já está instalado via npm install
		// Usar caminho absoluto do logo para evitar problemas
		absLogoPath, _ := filepath.Abs(logoPath)
		iconGenieCmd := exec.Command("npx", "icongenie", "generate", "-i", absLogoPath, "-m", "electron")
		iconGenieCmd.Stdout = os.Stdout
		iconGenieCmd.Stderr = os.Stderr
		
		if err := iconGenieCmd.Run(); err != nil {
			fmt.Println("⚠️  Aviso: Erro ao gerar ícones. Execute 'npm run generate:icons' manualmente.")
			fmt.Println("   Certifique-se de que @quasar/icongenie está instalado.")
		} else {
			fmt.Println("✓ Ícones do Electron gerados")
		}
	} else {
		fmt.Println("ℹ️  Logo não encontrado em src/assets/logo.png - ícones não serão gerados")
		fmt.Println("   Adicione o logo e execute 'npm run generate:icons' manualmente")
	}

	// Voltar para diretório original
	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("erro ao voltar para diretório original: %w", err)
	}

	return nil
}

// GenerateWebFrontend gera a estrutura frontend para Web (SPA)
func GenerateWebFrontend(projectName string, projectConfig *config.ProjectConfig) error {
	gen := templates.New(projectName)

	frontendConfig := FrontendConfig{
		ProjectName: projectName,
		ServerPort:  projectConfig.ServerPort,
	}

	// Criar estrutura de pastas frontend
	frontendDirs := []string{
		filepath.Join(projectName, "frontend", "src", "composables"),
		filepath.Join(projectName, "frontend", "src", "api"),
		filepath.Join(projectName, "frontend", "src", "components"),
		filepath.Join(projectName, "frontend", "src", "pages"),
		filepath.Join(projectName, "frontend", "src", "layouts"),
		filepath.Join(projectName, "frontend", "src", "router"),
		filepath.Join(projectName, "frontend", "src", "boot"),
		filepath.Join(projectName, "frontend", "src", "assets"),
		filepath.Join(projectName, "frontend", "src", "css"),
	}

	for _, dir := range frontendDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório %s: %w", dir, err)
		}
	}

	// Gerar arquivos frontend
	frontendFiles := map[string]string{
		"quasar_config.tmpl":            "frontend/quasar.config.js",
		"package_json_web.tmpl":         "frontend/package.json",
		"frontend_env.tmpl":             "frontend/.env",
		"composable_api.tmpl":           "frontend/src/composables/useApi.ts",
		"api_client.tmpl":               "frontend/src/api/client.js",
		"router_config.tmpl":            "frontend/src/router/index.js",
		"router_routes.tmpl":            "frontend/src/router/routes.js",
		"app_main.tmpl":                 "frontend/src/main.js",
		"app_vue.tmpl":                  "frontend/src/App.vue",
		"index_html.tmpl":               "frontend/index.html",
		"layout_main.tmpl":              "frontend/src/layouts/MainLayout.vue",
		"page_index.tmpl":               "frontend/src/pages/IndexPage.vue",
		"page_error.tmpl":               "frontend/src/pages/ErrorNotFound.vue",
		"component_essential_link.tmpl": "frontend/src/components/EssentialLink.vue",
		"boot_axios.tmpl":               "frontend/src/boot/axios.js",
		"app_scss.tmpl":                 "frontend/src/css/app.scss",
	}

	for template, output := range frontendFiles {
		if err := gen.Generate(template, output, frontendConfig); err != nil {
			return fmt.Errorf("erro ao gerar %s: %w", output, err)
		}
	}

	fmt.Println("✓ Estrutura frontend Web criada")

	// Instalar dependências npm
	fmt.Println("📦 Instalando dependências npm...")
	frontendPath := filepath.Join(projectName, "frontend")

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}

	// Mudar para diretório frontend
	if err := os.Chdir(frontendPath); err != nil {
		return fmt.Errorf("erro ao mudar para diretório frontend: %w", err)
	}

	npmCmd := exec.Command("npm", "install")
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr

	if err := npmCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao instalar dependências npm. Execute 'npm install' manualmente no diretório frontend.")
	} else {
		fmt.Println("✓ Dependências npm instaladas")
	}

	// Voltar para diretório original
	if err := os.Chdir(originalDir); err != nil {
		return fmt.Errorf("erro ao voltar para diretório original: %w", err)
	}

	return nil
}
