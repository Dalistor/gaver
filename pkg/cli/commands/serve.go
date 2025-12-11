package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/Dalistor/gaver/pkg/config"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Inicia o servidor de desenvolvimento",
		Long:  "Executa o servidor da aplicação em modo de desenvolvimento. Para projetos Mobile/Desktop, também inicia o Quasar dev server.",
		RunE:  runServe,
	}

	cmd.Flags().Bool("studio", false, "Abre Android Studio ou Xcode para debug (apenas para projetos Mobile)")
	cmd.Flags().String("platform", "", "Plataforma específica para Mobile (android ou ios)")
	cmd.Flags().Bool("cgo", false, "Habilita CGO para SQLite (requer compilador C). Se desabilitado, usa modernc.org/sqlite (puro Go)")

	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	// Verificar se estamos em um projeto Gaver
	if _, err := os.Stat("cmd/server/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("não parece ser um projeto Gaver. Execute 'gaver init' primeiro")
	}

	// Ler configuração do projeto
	projectConfig, err := config.ReadProjectConfig()
	if err != nil {
		// Se não encontrar GaverProject.json, assume tipo server
		enableCGO, _ := cmd.Flags().GetBool("cgo")
		return runNormal(enableCGO)
	}

	// Obter flag CGO
	enableCGO, _ := cmd.Flags().GetBool("cgo")

	// Executar baseado no tipo de projeto
	switch projectConfig.Type {
	case config.ProjectTypeMobile:
		openStudio, _ := cmd.Flags().GetBool("studio")
		platform, _ := cmd.Flags().GetString("platform")
		return runMobile(projectConfig, openStudio, platform, enableCGO)
	case config.ProjectTypeDesktop:
		return runDesktop(projectConfig, enableCGO)
	case config.ProjectTypeWeb:
		return runWeb(projectConfig, enableCGO)
	case config.ProjectTypeServer:
		return runNormal(enableCGO)
	default:
		return runNormal(enableCGO)
	}
}

func runNormal(enableCGO bool) error {
	fmt.Printf("Iniciando servidor...\n")
	fmt.Println("Use Ctrl+C para parar o servidor")

	// Verificar e configurar SQLite
	projectConfig, err := config.ReadProjectConfig()
	if err == nil && projectConfig.Database == "sqlite" {
		// Verificar se modernc.org/sqlite está disponível e funcionando
		hasModernc, err := checkSQLiteDependency()
		if err != nil || !hasModernc {
			// Se não tem modernc.org/sqlite ou erro ao adicionar, habilitar CGO
			if !enableCGO {
				if err != nil {
					fmt.Printf("⚠️  Aviso: %v\n", err)
				} else {
					fmt.Println("⚠️  modernc.org/sqlite não disponível")
				}
				fmt.Println("ℹ️  Tentando usar CGO com go-sqlite3...")
				enableCGO = true
			}

			// Se CGO foi habilitado, verificar se está disponível
			if enableCGO {
				if !checkCGOAvailable() {
					printCGOInstructions()
					return fmt.Errorf("CGO necessário mas não disponível. Siga as instruções acima")
				}
			}
		}
	}

	// Executar go run
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Configurar CGO para SQLite (projectConfig já foi lido acima)
	if projectConfig != nil && projectConfig.Database == "sqlite" {
		if !enableCGO {
			// Desabilitar CGO para forçar uso de modernc.org/sqlite
			cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
			fmt.Println("ℹ️  CGO desabilitado - usando modernc.org/sqlite (puro Go)")
		} else {
			// Habilitar CGO para usar go-sqlite3
			cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
			fmt.Println("ℹ️  CGO habilitado - usando go-sqlite3 (requer compilador C)")
		}
	}

	// Capturar sinais para parar gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Iniciar comando em goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Run()
	}()

	// Esperar por sinal ou erro
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("erro ao executar servidor: %w", err)
		}
	case <-sigChan:
		fmt.Println("\n\n🛑 Parando servidor...")
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}

	return nil
}

// checkCGOAvailable verifica se CGO está disponível no sistema
// Retorna true se CGO está disponível, false caso contrário
// Se não estiver disponível, imprime instruções para o usuário
func checkCGOAvailable() bool {
	// Tentar compilar um programa simples que usa CGO
	testCmd := exec.Command("go", "env", "CGO_ENABLED")
	testCmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	output, err := testCmd.Output()
	if err != nil {
		return false
	}

	// Verificar se CGO_ENABLED pode ser 1
	if strings.TrimSpace(string(output)) != "1" {
		return false
	}

	// Tentar verificar se há compilador C disponível
	// No Windows, verificar se gcc ou clang está disponível
	// No Linux/Mac, verificar se gcc ou clang está disponível
	if runtime.GOOS == "windows" {
		// Verificar MinGW ou TDM-GCC
		testGCC := exec.Command("gcc", "--version")
		if err := testGCC.Run(); err == nil {
			return true
		}
		// Verificar se está usando MSYS2 ou similar
		testClang := exec.Command("clang", "--version")
		if err := testClang.Run(); err == nil {
			return true
		}
	} else {
		// Linux/Mac - verificar gcc ou clang
		testGCC := exec.Command("gcc", "--version")
		if err := testGCC.Run(); err == nil {
			return true
		}
		testClang := exec.Command("clang", "--version")
		if err := testClang.Run(); err == nil {
			return true
		}
	}

	// Se chegou aqui, não há compilador C disponível
	return false
}

// printCGOInstructions imprime instruções sobre como habilitar CGO e instalar dependências
func printCGOInstructions() {
	fmt.Println("\n❌ CGO não está disponível no seu sistema.")
	fmt.Println("\n📋 Para usar SQLite com CGO, você precisa:")
	fmt.Println()

	if runtime.GOOS == "windows" {
		fmt.Println("1. Instalar um compilador C:")
		fmt.Println("   - Opção 1: MinGW-w64 (recomendado)")
		fmt.Println("     • Baixe em: https://www.mingw-w64.org/downloads/")
		fmt.Println("     • Ou instale via MSYS2: https://www.msys2.org/")
		fmt.Println("   - Opção 2: TDM-GCC")
		fmt.Println("     • Baixe em: https://jmeubank.github.io/tdm-gcc/")
		fmt.Println()
		fmt.Println("2. Adicionar o compilador ao PATH do sistema")
		fmt.Println()
		fmt.Println("3. Executar o comando com a flag --cgo:")
		fmt.Println("   gaver serve --cgo")
	} else if runtime.GOOS == "linux" {
		fmt.Println("1. Instalar compilador C e ferramentas de desenvolvimento:")
		fmt.Println("   • Ubuntu/Debian: sudo apt-get install build-essential")
		fmt.Println("   • Fedora: sudo dnf install gcc")
		fmt.Println("   • Arch: sudo pacman -S base-devel")
		fmt.Println()
		fmt.Println("2. Executar o comando com a flag --cgo:")
		fmt.Println("   gaver serve --cgo")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("1. Instalar Xcode Command Line Tools:")
		fmt.Println("   xcode-select --install")
		fmt.Println()
		fmt.Println("2. Executar o comando com a flag --cgo:")
		fmt.Println("   gaver serve --cgo")
	}

	fmt.Println()
	fmt.Println("💡 Alternativa: Use o driver SQLite puro Go (sem CGO)")
	fmt.Println("   O projeto já está configurado para usar github.com/glebarez/sqlite")
	fmt.Println("   que não requer CGO. Se você está vendo este erro, pode ser que")
	fmt.Println("   o driver puro Go não esteja disponível. Execute:")
	fmt.Println("   go get github.com/glebarez/sqlite")
	fmt.Println("   go mod tidy")
	fmt.Println()
}

// checkSQLiteDependency verifica se o projeto usa SQLite e se tem github.com/glebarez/sqlite
// Se não tiver, tenta adicionar automaticamente. Retorna true se glebarez/sqlite está disponível
func checkSQLiteDependency() (bool, error) {
	// Ler GaverProject.json para verificar tipo de banco
	projectConfig, err := config.ReadProjectConfig()
	if err != nil {
		return false, nil // Não é um erro crítico
	}

	if projectConfig.Database != "sqlite" {
		return false, nil // Não é SQLite
	}

	// Verificar se go.mod tem github.com/glebarez/sqlite
	goModPath := "go.mod"
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return false, nil
	}

	goModContent, err := os.ReadFile(goModPath)
	if err != nil {
		return false, nil
	}

	goModStr := string(goModContent)
	if strings.Contains(goModStr, "github.com/glebarez/sqlite") {
		// Verificar se database.go tem o import correto
		dbGoPath := filepath.Join("config", "database", "database.go")
		if dbContent, err := os.ReadFile(dbGoPath); err == nil {
			if strings.Contains(string(dbContent), "github.com/glebarez/sqlite") {
				return true, nil // Tudo OK, glebarez/sqlite disponível
			}
		}
		return true, nil // Tem no go.mod, mesmo que não esteja no código ainda
	}

	// Tentar adicionar automaticamente usando go get
	fmt.Println("⚠️  Tentando adicionar 'github.com/glebarez/sqlite' ao go.mod...")
	getCmd := exec.Command("go", "get", "github.com/glebarez/sqlite@v1.11.0")
	getCmd.Stdout = os.Stdout
	getCmd.Stderr = os.Stderr
	if err := getCmd.Run(); err != nil {
		// Se falhar, retornar false para habilitar CGO
		return false, fmt.Errorf("erro ao adicionar 'github.com/glebarez/sqlite': %w", err)
	}

	fmt.Println("✓ 'github.com/glebarez/sqlite' adicionado com sucesso")
	return true, nil
}

// waitForServer aguarda o servidor Go estar pronto fazendo requisições HTTP
func waitForServer(port string, maxAttempts int) error {
	url := fmt.Sprintf("http://localhost:%s", port)
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	for i := 0; i < maxAttempts; i++ {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			// Qualquer resposta HTTP significa que o servidor está rodando
			return nil
		}
		// Se não é erro de conexão, pode ser outro problema
		// Mas vamos continuar tentando
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("servidor não respondeu após %d tentativas", maxAttempts)
}

// readCapacitorAppID lê o app ID do capacitor.config.json ou capacitor.config.js
func readCapacitorAppID(frontendPath string) string {
	// Tentar capacitor.config.json primeiro
	configPath := filepath.Join(frontendPath, "capacitor.config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var config struct {
			AppID string `json:"appId"`
		}
		if err := json.Unmarshal(data, &config); err == nil && config.AppID != "" {
			return config.AppID
		}
	}

	// Tentar capacitor.config.js como fallback
	configPathJS := filepath.Join(frontendPath, "capacitor.config.js")
	if data, err := os.ReadFile(configPathJS); err == nil {
		content := string(data)
		// Procurar por appId: "..." ou appId: '...'
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.Contains(line, "appId") {
				// Extrair valor entre aspas
				if idx := strings.Index(line, `"`); idx != -1 {
					start := idx + 1
					if end := strings.Index(line[start:], `"`); end != -1 {
						appID := line[start : start+end]
						if appID != "" {
							return appID
						}
					}
				}
			}
		}
	}

	return ""
}

func runMobile(projectConfig *config.ProjectConfig, openStudio bool, platform string, enableCGO bool) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

	// Determinar plataforma (android ou ios)
	if platform == "" {
		platform = "android" // Padrão
	}
	if platform != "android" && platform != "ios" {
		return fmt.Errorf("plataforma inválida: %s. Use 'android' ou 'ios'", platform)
	}

	// Verificar e configurar SQLite
	projectConfig, err := config.ReadProjectConfig()
	if err == nil && projectConfig.Database == "sqlite" {
		// Verificar se modernc.org/sqlite está disponível e funcionando
		hasModernc, err := checkSQLiteDependency()
		if err != nil || !hasModernc {
			// Se não tem modernc.org/sqlite ou erro ao adicionar, habilitar CGO
			if !enableCGO {
				if err != nil {
					fmt.Printf("⚠️  Aviso: %v\n", err)
				} else {
					fmt.Println("⚠️  modernc.org/sqlite não disponível")
				}
				fmt.Println("ℹ️  Tentando usar CGO com go-sqlite3...")
				enableCGO = true
			}

			// Se CGO foi habilitado, verificar se está disponível
			if enableCGO {
				if !checkCGOAvailable() {
					printCGOInstructions()
					return fmt.Errorf("CGO necessário mas não disponível. Siga as instruções acima")
				}
			}
		}
	}

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	// Iniciar servidor Go em goroutine
	goServer := exec.Command("go", "run", "cmd/server/main.go")
	goServer.Stdout = os.Stdout
	goServer.Stderr = os.Stderr

	// Configurar CGO para SQLite
	if projectConfig.Database == "sqlite" {
		if !enableCGO {
			// Desabilitar CGO para forçar uso de modernc.org/sqlite
			goServer.Env = append(os.Environ(), "CGO_ENABLED=0")
			fmt.Println("ℹ️  CGO desabilitado - usando modernc.org/sqlite (puro Go)")
		} else {
			// Habilitar CGO para usar go-sqlite3
			goServer.Env = append(os.Environ(), "CGO_ENABLED=1")
			fmt.Println("ℹ️  CGO habilitado - usando go-sqlite3 (requer compilador C)")
		}
	}

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Aguardar servidor Go estar pronto antes de iniciar Quasar
	serverPort := projectConfig.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}
	fmt.Printf("⏳ Aguardando servidor Go na porta %s...\n", serverPort)
	if err := waitForServer(serverPort, 30); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao aguardar servidor Go: %w", err)
	}
	fmt.Println("✓ Servidor Go está pronto")

	// Verificar se o diretório android existe antes de iniciar Quasar
	if platform == "android" {
		androidPath := filepath.Join(frontendPath, "android")
		if _, err := os.Stat(androidPath); os.IsNotExist(err) {
			goServer.Process.Kill()
			return fmt.Errorf("diretório android não encontrado. Execute 'npx cap add android' primeiro no diretório frontend")
		}
	}

	// Tentar ler app ID do capacitor.config.json para passar automaticamente
	appID := readCapacitorAppID(frontendPath)
	if appID != "" {
		fmt.Printf("📱 App ID detectado: %s\n", appID)
		fmt.Println("ℹ️  O Quasar será configurado automaticamente com este app ID")
	} else {
		fmt.Println("ℹ️  App ID não encontrado no capacitor.config.json")
		fmt.Println("ℹ️  O Quasar pode perguntar sobre o app ID - pressione Enter para usar o padrão")
	}

	// Iniciar Quasar dev server
	fmt.Printf("🚀 Iniciando Quasar dev server para %s...\n", platform)
	quasarCmd := exec.Command("npx", "quasar", "dev", "-m", "capacitor", "-T", platform)
	quasarCmd.Dir = frontendPath
	quasarCmd.Stdout = os.Stdout
	quasarCmd.Stderr = os.Stderr

	// Criar pipe para stdin se tivermos app ID
	var stdinPipe io.WriteCloser
	if appID != "" {
		var err error
		stdinPipe, err = quasarCmd.StdinPipe()
		if err == nil {
			quasarCmd.Stdin = nil // Não usar os.Stdin quando temos pipe
		} else {
			// Se falhar ao criar pipe, usar stdin normal
			quasarCmd.Stdin = os.Stdin
			stdinPipe = nil
		}
	} else {
		quasarCmd.Stdin = os.Stdin // Permitir entrada interativa
	}

	fmt.Printf("📂 Diretório: %s\n", frontendPath)
	fmt.Printf("🔧 Comando: npx quasar dev -m capacitor -T %s\n", platform)
	fmt.Println("ℹ️  Aguarde o Quasar compilar e iniciar...")

	// Iniciar Quasar em goroutine para não bloquear
	quasarErrChan := make(chan error, 1)
	go func() {
		if stdinPipe != nil {
			// Aguardar mais tempo para o Quasar fazer a pergunta
			// O Quasar geralmente pergunta após alguns segundos
			time.Sleep(3 * time.Second)
			// Enviar o app ID automaticamente
			io.WriteString(stdinPipe, appID+"\n")
			stdinPipe.Close()
		}
		quasarErrChan <- quasarCmd.Run()
	}()

	// Aguardar um pouco para verificar se há erro imediato
	time.Sleep(500 * time.Millisecond)
	select {
	case err := <-quasarErrChan:
		// Quasar terminou muito rápido, provavelmente erro
		goServer.Process.Kill()
		return fmt.Errorf("Quasar dev server terminou com erro: %w", err)
	default:
		// Quasar ainda está rodando, tudo OK
		fmt.Println("✓ Quasar dev server iniciado")
		fmt.Println("ℹ️  O Quasar pode levar alguns segundos para compilar pela primeira vez")
	}

	// Se flag --studio, abrir Android Studio ou Xcode
	if openStudio {
		if platform == "android" {
			androidPath := filepath.Join(frontendPath, "android")
			if _, err := os.Stat(androidPath); err == nil {
				fmt.Println("📱 Abrindo Android Studio...")
				studioCmd := exec.Command("studio", androidPath)
				studioCmd.Start() // Não esperar, apenas iniciar
			}
		} else if platform == "ios" {
			iosPath := filepath.Join(frontendPath, "ios")
			if _, err := os.Stat(iosPath); err == nil {
				fmt.Println("🍎 Abrindo Xcode...")
				// No macOS, usar 'open' para abrir Xcode
				if runtime.GOOS == "darwin" {
					studioCmd := exec.Command("open", "-a", "Xcode", iosPath)
					studioCmd.Start()
				} else {
					fmt.Println("⚠️  Xcode só está disponível no macOS")
				}
			}
		}
	}

	// Capturar sinais para parar gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Esperar por sinal ou erro do Quasar
	select {
	case <-sigChan:
		fmt.Println("\n\n🛑 Parando servidores...")
		if goServer.Process != nil {
			goServer.Process.Kill()
		}
		if quasarCmd.Process != nil {
			quasarCmd.Process.Kill()
		}
		return nil
	case err := <-quasarErrChan:
		// Quasar terminou (pode ser erro ou término normal)
		if err != nil {
			fmt.Printf("\n⚠️  Quasar dev server terminou: %v\n", err)
		}
		if goServer.Process != nil {
			goServer.Process.Kill()
		}
		return err
	}
}

func runDesktop(projectConfig *config.ProjectConfig, enableCGO bool) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

	// Verificar e configurar SQLite
	projectConfig, err := config.ReadProjectConfig()
	if err == nil && projectConfig.Database == "sqlite" {
		// Verificar se modernc.org/sqlite está disponível e funcionando
		hasModernc, err := checkSQLiteDependency()
		if err != nil || !hasModernc {
			// Se não tem modernc.org/sqlite ou erro ao adicionar, habilitar CGO
			if !enableCGO {
				if err != nil {
					fmt.Printf("⚠️  Aviso: %v\n", err)
				} else {
					fmt.Println("⚠️  modernc.org/sqlite não disponível")
				}
				fmt.Println("ℹ️  Habilitando CGO automaticamente para usar go-sqlite3")
				enableCGO = true
			}
		}
	}

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	// Iniciar servidor Go em goroutine
	goServer := exec.Command("go", "run", "cmd/server/main.go")
	goServer.Stdout = os.Stdout
	goServer.Stderr = os.Stderr

	// Configurar CGO para SQLite
	if projectConfig.Database == "sqlite" {
		if !enableCGO {
			// Desabilitar CGO para forçar uso de modernc.org/sqlite
			goServer.Env = append(os.Environ(), "CGO_ENABLED=0")
		} else {
			// Habilitar CGO para usar go-sqlite3
			goServer.Env = append(os.Environ(), "CGO_ENABLED=1")
		}
	}

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Aguardar servidor Go estar pronto antes de iniciar Electron
	serverPort := projectConfig.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}
	fmt.Printf("⏳ Aguardando servidor Go na porta %s...\n", serverPort)
	if err := waitForServer(serverPort, 30); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao aguardar servidor Go: %w", err)
	}
	fmt.Println("✓ Servidor Go está pronto")

	// Iniciar Quasar dev server
	quasarCmd := exec.Command("npx", "quasar", "dev", "-m", "electron")
	quasarCmd.Dir = frontendPath
	quasarCmd.Stdout = os.Stdout
	quasarCmd.Stderr = os.Stderr

	if err := quasarCmd.Start(); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao iniciar Quasar dev server: %w", err)
	}

	// Capturar sinais para parar gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Esperar por sinal
	<-sigChan
	fmt.Println("\n\n🛑 Parando servidores...")

	if goServer.Process != nil {
		goServer.Process.Kill()
	}
	if quasarCmd.Process != nil {
		quasarCmd.Process.Kill()
	}

	return nil
}

func runWeb(projectConfig *config.ProjectConfig, enableCGO bool) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

	// Verificar e configurar SQLite
	projectConfig, err := config.ReadProjectConfig()
	if err == nil && projectConfig.Database == "sqlite" {
		// Verificar se modernc.org/sqlite está disponível e funcionando
		hasModernc, err := checkSQLiteDependency()
		if err != nil || !hasModernc {
			// Se não tem modernc.org/sqlite ou erro ao adicionar, habilitar CGO
			if !enableCGO {
				if err != nil {
					fmt.Printf("⚠️  Aviso: %v\n", err)
				} else {
					fmt.Println("⚠️  modernc.org/sqlite não disponível")
				}
				fmt.Println("ℹ️  Tentando usar CGO com go-sqlite3...")
				enableCGO = true
			}

			// Se CGO foi habilitado, verificar se está disponível
			if enableCGO {
				if !checkCGOAvailable() {
					printCGOInstructions()
					return fmt.Errorf("CGO necessário mas não disponível. Siga as instruções acima")
				}
			}
		}
	}

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	// Iniciar servidor Go em goroutine
	goServer := exec.Command("go", "run", "cmd/server/main.go")
	goServer.Stdout = os.Stdout
	goServer.Stderr = os.Stderr

	// Configurar CGO para SQLite
	if projectConfig.Database == "sqlite" {
		if !enableCGO {
			// Desabilitar CGO para forçar uso de modernc.org/sqlite
			goServer.Env = append(os.Environ(), "CGO_ENABLED=0")
			fmt.Println("ℹ️  CGO desabilitado - usando modernc.org/sqlite (puro Go)")
		} else {
			// Habilitar CGO para usar go-sqlite3
			goServer.Env = append(os.Environ(), "CGO_ENABLED=1")
			fmt.Println("ℹ️  CGO habilitado - usando go-sqlite3 (requer compilador C)")
		}
	}

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Aguardar servidor Go estar pronto antes de iniciar Quasar
	serverPort := projectConfig.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}
	fmt.Printf("⏳ Aguardando servidor Go na porta %s...\n", serverPort)
	if err := waitForServer(serverPort, 30); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao aguardar servidor Go: %w", err)
	}
	fmt.Println("✓ Servidor Go está pronto")

	// Iniciar Quasar dev server
	quasarCmd := exec.Command("npx", "quasar", "dev")
	quasarCmd.Dir = frontendPath
	quasarCmd.Stdout = os.Stdout
	quasarCmd.Stderr = os.Stderr

	if err := quasarCmd.Start(); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao iniciar Quasar dev server: %w", err)
	}

	// Capturar sinais para parar gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Esperar por sinal
	<-sigChan
	fmt.Println("\n\n🛑 Parando servidores...")

	if goServer.Process != nil {
		goServer.Process.Kill()
	}
	if quasarCmd.Process != nil {
		quasarCmd.Process.Kill()
	}

	return nil
}
