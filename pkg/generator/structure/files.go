package structure

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

// sanitizeForAppId sanitiza o nome do projeto para uso em appId (package name)
func sanitizeForAppId(name string) string {
	// Converter para minúsculas e remover caracteres inválidos
	result := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result += string(r)
		} else if r >= 'A' && r <= 'Z' {
			result += string(r + 32) // Converter para minúscula
		} else if r == '-' || r == '_' {
			result += string(r)
		} else if r == ' ' {
			result += "" // Remover espaços
		}
		// Ignorar outros caracteres
	}
	if result == "" {
		result = "app"
	}
	return result
}

// min retorna o menor valor entre dois inteiros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GenerateMobileFrontend gera a estrutura frontend para Mobile (Android + iOS)
func GenerateMobileFrontend(projectName string, projectConfig *config.ProjectConfig) error {
	gen := templates.New(projectName)

	// Sanitizar nome para appId
	sanitizedAppId := sanitizeForAppId(projectName)
	if sanitizedAppId == "" {
		sanitizedAppId = "app" // Fallback se sanitização resultar em string vazia
	}

	frontendConfig := FrontendConfig{
		ProjectName: sanitizedAppId, // Usar versão sanitizada para appId
		ServerPort:  projectConfig.ServerPort,
	}

	// Debug: garantir que ProjectName não está vazio
	if frontendConfig.ProjectName == "" {
		frontendConfig.ProjectName = "app"
		fmt.Println("⚠️  Aviso: ProjectName estava vazio, usando 'app' como padrão")
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

	// Gerar arquivos frontend (sem capacitor.config.js - será gerado após cap init)
	frontendFiles := map[string]string{
		"quasar_config.tmpl":            "frontend/quasar.config.js",
		"package_json_mobile.tmpl":      "frontend/package.json",
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

	// Remover capacitor.config.js/.ts/.json se existir (pode ter sido gerado incorretamente antes)
	capacitorConfigPathJs := "capacitor.config.js"
	capacitorConfigPathTs := "capacitor.config.ts"
	capacitorConfigPathJson := "capacitor.config.json"
	if _, err := os.Stat(capacitorConfigPathJs); err == nil {
		os.Remove(capacitorConfigPathJs)
	}
	if _, err := os.Stat(capacitorConfigPathTs); err == nil {
		os.Remove(capacitorConfigPathTs)
	}
	if _, err := os.Stat(capacitorConfigPathJson); err == nil {
		os.Remove(capacitorConfigPathJson)
	}

	// Gerar capacitor.config.json diretamente (não usar cap init para evitar problemas de modo não-interativo)
	fmt.Println("🔧 Configurando capacitor.config.json...")

	// Garantir que ProjectName não está vazio
	appIdName := frontendConfig.ProjectName
	if appIdName == "" {
		appIdName = "app"
		fmt.Println("⚠️  ProjectName estava vazio, usando 'app' como padrão")
	}

	// Usar projectName original para appName (não sanitizado)
	appDisplayName := projectName
	if appDisplayName == "" {
		appDisplayName = "Gaver App"
	}

	// Escapar caracteres especiais no appDisplayName para evitar problemas
	// Escapar aspas duplas, quebras de linha e outros caracteres problemáticos
	appDisplayNameEscaped := appDisplayName
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "\\", "\\\\") // Escapar backslashes primeiro
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "\"", "\\\"") // Escapar aspas duplas
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "'", "\\'")   // Escapar aspas simples
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "\n", " ")    // Remover quebras de linha
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "\r", " ")    // Remover carriage return
	appDisplayNameEscaped = strings.ReplaceAll(appDisplayNameEscaped, "\t", " ")    // Remover tabs

	// Garantir que appIdName também está limpo
	appIdNameClean := strings.TrimSpace(appIdName)
	if appIdNameClean == "" {
		appIdNameClean = "app"
	}

	// Gerar conteúdo JSON diretamente - mais simples e não tem problemas com ESM/CommonJS
	// JSON é o formato mais confiável e não depende do tipo de módulo
	configContent := fmt.Sprintf(`{
  "appId": "com.%s.app",
  "appName": "%s",
  "webDir": "dist",
  "server": {
    "androidScheme": "https",
    "iosScheme": "https"
  },
  "plugins": {
    "Filesystem": {
      "android": {
        "path": "files"
      },
      "ios": {
        "path": "files"
      }
    }
  }
}
`, appIdNameClean, appDisplayNameEscaped)

	// Escrever arquivo diretamente no diretório atual (frontend)
	// Usar .json que é mais simples e não tem problemas com ESM/CommonJS
	// IMPORTANTE: garantir que estamos no diretório frontend
	configPath := "capacitor.config.json"

	// Remover .js/.ts se existir (migrar para .json)
	oldJsPath := "capacitor.config.js"
	oldTsPath := "capacitor.config.ts"
	if _, err := os.Stat(oldJsPath); err == nil {
		os.Remove(oldJsPath)
		fmt.Println("   Removido capacitor.config.js antigo (migrando para .json)")
	}
	if _, err := os.Stat(oldTsPath); err == nil {
		os.Remove(oldTsPath)
		fmt.Println("   Removido capacitor.config.ts antigo (migrando para .json)")
	}

	// Verificar diretório atual para debug
	if currentDir, err := os.Getwd(); err == nil {
		fmt.Printf("   Diretório atual: %s\n", currentDir)
	}

	// Debug: mostrar conteúdo que será escrito
	fmt.Printf("   Gerando arquivo com appId: com.%s.app\n", appIdNameClean)
	fmt.Printf("   appName: %s\n", appDisplayNameEscaped)

	// Validar que o conteúdo não tem problemas óbvios
	if strings.Contains(configContent, "{{") || strings.Contains(configContent, "}}") {
		fmt.Println("⚠️  ERRO CRÍTICO: Conteúdo contém variáveis de template não processadas!")
		fmt.Printf("   Conteúdo: %s\n", configContent)
		return fmt.Errorf("erro ao gerar capacitor.config.json: variáveis não processadas")
	}

	// Verificar se há problemas de sintaxe básicos (JSON deve ter appId)
	if !strings.Contains(configContent, `"appId"`) {
		fmt.Println("⚠️  ERRO: Conteúdo não contém appId!")
		fmt.Printf("   Conteúdo: %s\n", configContent)
		return fmt.Errorf("erro ao gerar capacitor.config.json: appId ausente")
	}

	// Escrever arquivo com encoding UTF-8 explícito
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao escrever capacitor.config.json: %v\n", err)
		return fmt.Errorf("erro ao escrever capacitor.config.json: %w", err)
	}

	// Verificar se o arquivo foi criado e ler de volta para validar
	if fileInfo, err := os.Stat(configPath); err != nil {
		fmt.Printf("⚠️  Aviso: Arquivo não encontrado após escrita: %v\n", err)
		return fmt.Errorf("arquivo não encontrado após escrita: %w", err)
	} else {
		fmt.Printf("   Arquivo criado: %d bytes\n", fileInfo.Size())

		// Ler arquivo de volta para validar
		readBack, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao ler arquivo de volta: %v\n", err)
			return fmt.Errorf("erro ao ler arquivo de volta: %w", err)
		}

		readBackStr := string(readBack)

		// Verificações rigorosas
		hasErrors := false
		// Validar JSON
		if strings.Contains(readBackStr, "{{") || strings.Contains(readBackStr, "}}") {
			fmt.Println("❌ ERRO CRÍTICO: capacitor.config.json contém variáveis não processadas!")
			fmt.Printf("   Conteúdo completo:\n%s\n", readBackStr)
			hasErrors = true
		}

		if strings.Contains(readBackStr, `"appId": ""`) || strings.Contains(readBackStr, `"appId": "com..app"`) {
			fmt.Println("❌ ERRO CRÍTICO: appId está vazio ou inválido!")
			fmt.Printf("   Conteúdo: %s\n", readBackStr)
			hasErrors = true
		}

		if !strings.Contains(readBackStr, `"appId"`) {
			fmt.Println("❌ ERRO CRÍTICO: capacitor.config.json não contém appId!")
			fmt.Printf("   Conteúdo: %s\n", readBackStr)
			hasErrors = true
		}

		// Validar que é JSON válido (deve começar com { e terminar com })
		if !strings.HasPrefix(strings.TrimSpace(readBackStr), "{") || !strings.HasSuffix(strings.TrimSpace(readBackStr), "}") {
			fmt.Println("❌ ERRO CRÍTICO: capacitor.config.json não é JSON válido!")
			fmt.Printf("   Conteúdo: %s\n", readBackStr)
			hasErrors = true
		}

		if hasErrors {
			return fmt.Errorf("capacitor.config.json gerado com erros - verifique o conteúdo acima")
		}

		fmt.Println("✓ capacitor.config.json configurado e validado")
		fmt.Printf("   appId: com.%s.app\n", appIdNameClean)

		// Mostrar primeiras linhas para confirmação
		lines := strings.Split(readBackStr, "\n")
		if len(lines) > 3 {
			fmt.Printf("   Verificação: %s\n", strings.TrimSpace(lines[3]))
		}
	}

	// Não testar com Node.js pois o arquivo usa ESM (import/export)
	// O Capacitor vai validar quando tentar usar

	// Adicionar plataforma Android
	fmt.Println("📱 Adicionando plataforma Android...")
	capacitorAddAndroidCmd := exec.Command("npx", "cap", "add", "android")
	capacitorAddAndroidCmd.Stdout = os.Stdout
	capacitorAddAndroidCmd.Stderr = os.Stderr
	if err := capacitorAddAndroidCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao adicionar Android. Execute 'npx cap add android' manualmente se necessário.")
		fmt.Println("   Possíveis causas:")
		fmt.Println("   - capacitor.config.json tem sintaxe inválida")
		fmt.Println("   - Capacitor não está instalado corretamente")
		fmt.Println("   - Verifique o arquivo capacitor.config.json manualmente")
	}

	// Adicionar plataforma iOS (apenas no macOS)
	if runtime.GOOS == "darwin" {
		fmt.Println("🍎 Adicionando plataforma iOS...")
		capacitorAddIOSCmd := exec.Command("npx", "cap", "add", "ios")
		capacitorAddIOSCmd.Stdout = os.Stdout
		capacitorAddIOSCmd.Stderr = os.Stderr
		if err := capacitorAddIOSCmd.Run(); err != nil {
			fmt.Println("⚠️  Aviso: Erro ao adicionar iOS. Execute 'npx cap add ios' manualmente se necessário.")
			fmt.Println("   Certifique-se de que o Xcode está instalado.")
		} else {
			fmt.Println("✓ Plataforma iOS adicionada")
		}
	} else {
		fmt.Println("ℹ️  iOS não será adicionado automaticamente")
		fmt.Println("   iOS requer macOS e Xcode instalado.")
		fmt.Println("   Para compilar para iOS, você precisa:")
		fmt.Println("   1. Usar um computador macOS")
		fmt.Println("   2. Ter o Xcode instalado")
		fmt.Println("   3. Executar: cd frontend && npx cap add ios")
		fmt.Println("   4. Depois: gaver build --platform ios")
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
		filepath.Join("assets", "logo.png"),             // Diretório atual (desenvolvimento)
		filepath.Join("..", "assets", "logo.png"),       // Um nível acima
		filepath.Join("..", "..", "assets", "logo.png"), // Dois níveis acima
		filepath.Join(currentDir, "assets", "logo.png"), // Diretório atual absoluto
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
