package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	templates "github.com/Dalistor/gaver/internal/templates"
	"github.com/Dalistor/gaver/pkg/config"
	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Compila o projeto para produção",
		Long:  "Gera os arquivos de build do projeto. Para Mobile gera APK (Android) ou IPA (iOS), para Desktop gera .exe, para Server faz build Go normal.",
		RunE:  runBuild,
	}

	cmd.Flags().String("platform", "", "Plataforma específica para Mobile (android ou ios). Se não especificado, compila ambas.")
	cmd.Flags().Bool("release", false, "Gera build de release assinado (apenas Android)")

	return cmd
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Verificar se estamos em um projeto Gaver
	if _, err := os.Stat("cmd/server/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("não parece ser um projeto Gaver. Execute 'gaver init' primeiro")
	}

	// Ler configuração do projeto
	projectConfig, err := config.ReadProjectConfig()
	if err != nil {
		// Se não encontrar GaverProject.json, assume tipo server
		return buildServer()
	}

	// Obter flags
	platform, _ := cmd.Flags().GetString("platform")
	release, _ := cmd.Flags().GetBool("release")

	// Executar build baseado no tipo de projeto
	switch projectConfig.Type {
	case config.ProjectTypeMobile:
		if platform == "" || platform == "android" {
			if err := buildMobile(projectConfig, "android", release); err != nil {
				return err
			}
		}
		if platform == "" || platform == "ios" {
			if err := buildIOS(projectConfig); err != nil {
				return err
			}
		}
		return nil
	case config.ProjectTypeDesktop:
		return buildDesktop(projectConfig)
	case config.ProjectTypeWeb:
		return buildWeb(projectConfig)
	case config.ProjectTypeServer:
		return buildServer()
	default:
		return buildServer()
	}
}

func buildServer() error {
	fmt.Println("🔨 Compilando servidor Go...")

	// Build do servidor Go
	buildCmd := exec.Command("go", "build", "-o", "bin/server", "cmd/server/main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar servidor: %w", err)
	}

	fmt.Println("✓ Build concluído! Executável em: bin/server")
	return nil
}

func buildMobile(projectConfig *config.ProjectConfig, platform string, release bool) error {
	if platform != "android" {
		return fmt.Errorf("plataforma não suportada: %s", platform)
	}

	fmt.Println("🔨 Compilando projeto Android...")

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	androidPath := filepath.Join(frontendPath, "android")

	// 1. Compilar servidor Go para Android (binário ARM64)
	fmt.Println("📱 Compilando servidor Go para Android...")
	serverBinaryPath := filepath.Join(androidPath, "app", "src", "main", "assets", "server")
	if err := os.MkdirAll(filepath.Dir(serverBinaryPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório assets: %w", err)
	}

	// Compilar para Android ARM64 (sempre sem CGO para Android)
	goBuildCmd := exec.Command("go", "build", "-o", serverBinaryPath, "-ldflags", "-s -w", "cmd/server/main.go")
	goBuildCmd.Env = append(os.Environ(), "GOOS=android", "GOARCH=arm64", "CGO_ENABLED=0")
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr

	var serverCompiled bool
	if err := goBuildCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao compilar servidor Go para Android. Continuando sem servidor...")
		serverCompiled = false
	} else {
		fmt.Println("✓ Servidor Go compilado para Android")
		serverCompiled = true
	}

	// Se usar SQLite, copiar banco de dados para assets (se existir)
	if projectConfig.Database == "sqlite" {
		dbName := projectConfig.ProjectName
		if dbName == "" {
			dbName = "app"
		}
		dbFile := dbName + ".db"

		// Procurar arquivo .db no diretório atual ou em subdiretórios comuns
		var dbPath string
		possiblePaths := []string{
			dbFile,
			filepath.Join("data", dbFile),
			filepath.Join("database", dbFile),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				dbPath = path
				break
			}
		}

		if dbPath != "" {
			fmt.Printf("📦 Copiando banco SQLite embutido: %s\n", dbPath)
			dbDest := filepath.Join(androidPath, "app", "src", "main", "assets", "database.db")
			if err := copyFile(dbPath, dbDest); err != nil {
				fmt.Printf("⚠️  Aviso: Erro ao copiar banco SQLite: %v\n", err)
			} else {
				fmt.Println("✓ Banco SQLite copiado para assets")
			}
		} else {
			fmt.Println("ℹ️  Nenhum banco SQLite encontrado - será criado em tempo de execução")
		}
	}

	// 2. Gerar AAR do Go usando gomobile (opcional, para funções Go)
	fmt.Println("📱 Gerando AAR do Go (opcional)...")
	if _, err := exec.LookPath("gomobile"); err == nil {
		// Criar diretório para o AAR se não existir
		libsPath := filepath.Join(androidPath, "app", "libs")
		if err := os.MkdirAll(libsPath, 0755); err != nil {
			fmt.Println("⚠️  Aviso: Erro ao criar diretório libs")
		} else {
			// Verificar se existe cmd/mobile, se não, criar estrutura básica
			if _, err := os.Stat("cmd/mobile"); os.IsNotExist(err) {
				if err := createMobilePackage(); err == nil {
					// Gerar AAR usando gomobile bind
					gomobileCmd := exec.Command("gomobile", "bind", "-target=android", "-o", filepath.Join(libsPath, "gaver.aar"), "./cmd/mobile")
					gomobileCmd.Stdout = os.Stdout
					gomobileCmd.Stderr = os.Stderr
					if err := gomobileCmd.Run(); err == nil {
						fmt.Println("✓ AAR gerado com sucesso")
					}
				}
			}
		}
	}

	// 2. Build do Quasar/Capacitor
	fmt.Println("📦 Compilando Quasar com Capacitor...")
	quasarBuildCmd := exec.Command("quasar", "build", "-m", "capacitor", "-T", "android")
	quasarBuildCmd.Dir = frontendPath
	quasarBuildCmd.Stdout = os.Stdout
	quasarBuildCmd.Stderr = os.Stderr

	if err := quasarBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar Quasar: %w", err)
	}

	// 3. Configurar build.gradle para incluir o AAR (se foi gerado)
	libsPath := filepath.Join(androidPath, "app", "libs")
	if _, err := os.Stat(filepath.Join(libsPath, "gaver.aar")); err == nil {
		fmt.Println("🔧 Configurando build.gradle para incluir AAR...")
		buildGradlePath := filepath.Join(androidPath, "app", "build.gradle")
		if err := configureGradleForAAR(buildGradlePath, libsPath); err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao configurar build.gradle: %v\n", err)
			fmt.Println("   Você precisará adicionar manualmente a dependência do AAR no build.gradle")
		} else {
			fmt.Println("✓ build.gradle configurado para incluir AAR")
		}
	}

	// 4. Gerar MainActivity.java que inicia o servidor (se servidor foi compilado)
	if serverCompiled {
		fmt.Println("🔧 Configurando MainActivity para iniciar servidor...")
		if err := generateMainActivity(androidPath, projectConfig); err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao gerar MainActivity: %v\n", err)
			fmt.Println("   Você precisará configurar manualmente o MainActivity para iniciar o servidor")
		} else {
			fmt.Println("✓ MainActivity configurado")
		}
	}

	// 4. Gerar/Verificar keystore para release (se necessário)
	if release {
		if err := generateAndroidKeystore(androidPath, projectConfig); err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao gerar/configurar keystore: %v\n", err)
			fmt.Println("   Continuando com build debug...")
			release = false
		}
	}

	// 5. Build do Android (gerar APK)
	fmt.Println("📱 Compilando projeto Android...")

	// Verificar se gradlew existe
	gradlewPath := filepath.Join(androidPath, "gradlew")
	if runtime.GOOS == "windows" {
		gradlewPath = filepath.Join(androidPath, "gradlew.bat")
	}

	if _, err := os.Stat(gradlewPath); os.IsNotExist(err) {
		return fmt.Errorf("gradlew não encontrado. Execute 'quasar capacitor sync android' primeiro")
	}

	// Escolher tipo de build
	buildType := "assembleDebug"
	if release {
		buildType = "assembleRelease"
	}

	gradleCmd := exec.Command(gradlewPath, buildType)
	gradleCmd.Dir = androidPath
	gradleCmd.Stdout = os.Stdout
	gradleCmd.Stderr = os.Stderr

	if err := gradleCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar Android: %w", err)
	}

	// Caminho do APK gerado
	apkType := "debug"
	if release {
		apkType = "release"
	}
	apkPath := filepath.Join(androidPath, "app", "build", "outputs", "apk", apkType, fmt.Sprintf("app-%s.apk", apkType))

	fmt.Printf("✓ Build concluído!\n")
	fmt.Printf("📱 APK: %s\n", apkPath)
	if _, err := os.Stat(filepath.Join(libsPath, "gaver.aar")); err == nil {
		fmt.Printf("📦 AAR incluído: %s\n", filepath.Join(libsPath, "gaver.aar"))
	}

	return nil
}

func buildIOS(projectConfig *config.ProjectConfig) error {
	fmt.Println("🔨 Compilando projeto iOS...")

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	iosPath := filepath.Join(frontendPath, "ios")

	// Verificar se estamos no macOS
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("build iOS requer macOS e Xcode instalado")
	}

	// 1. Compilar servidor Go para iOS (binário ARM64)
	fmt.Println("🍎 Compilando servidor Go para iOS...")
	serverBinaryPath := filepath.Join(iosPath, "App", "server")
	if err := os.MkdirAll(filepath.Dir(serverBinaryPath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	// Compilar para iOS ARM64 (sempre sem CGO para iOS)
	goBuildCmd := exec.Command("go", "build", "-o", serverBinaryPath, "-ldflags", "-s -w", "cmd/server/main.go")
	goBuildCmd.Env = append(os.Environ(), "GOOS=ios", "GOARCH=arm64", "CGO_ENABLED=0")
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr

	if err := goBuildCmd.Run(); err != nil {
		fmt.Println("⚠️  Aviso: Erro ao compilar servidor Go para iOS. Continuando sem servidor...")
	} else {
		fmt.Println("✓ Servidor Go compilado para iOS")
	}

	// Se usar SQLite, copiar banco de dados (se existir)
	if projectConfig.Database == "sqlite" {
		dbName := projectConfig.ProjectName
		if dbName == "" {
			dbName = "app"
		}
		dbFile := dbName + ".db"

		// Procurar arquivo .db no diretório atual ou em subdiretórios comuns
		var dbPath string
		possiblePaths := []string{
			dbFile,
			filepath.Join("data", dbFile),
			filepath.Join("database", dbFile),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				dbPath = path
				break
			}
		}

		if dbPath != "" {
			fmt.Printf("📦 Copiando banco SQLite embutido: %s\n", dbPath)
			dbDest := filepath.Join(iosPath, "App", "database.db")
			if err := copyFile(dbPath, dbDest); err != nil {
				fmt.Printf("⚠️  Aviso: Erro ao copiar banco SQLite: %v\n", err)
			} else {
				fmt.Println("✓ Banco SQLite copiado")
			}
		} else {
			fmt.Println("ℹ️  Nenhum banco SQLite encontrado - será criado em tempo de execução")
		}
	}

	// 2. Build do Quasar/Capacitor
	fmt.Println("📦 Compilando Quasar com Capacitor...")
	quasarBuildCmd := exec.Command("quasar", "build", "-m", "capacitor", "-T", "ios")
	quasarBuildCmd.Dir = frontendPath
	quasarBuildCmd.Stdout = os.Stdout
	quasarBuildCmd.Stderr = os.Stderr

	if err := quasarBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar Quasar: %w", err)
	}

	// 3. Sincronizar com Capacitor
	fmt.Println("🔄 Sincronizando com Capacitor...")
	capacitorSyncCmd := exec.Command("npx", "cap", "sync", "ios")
	capacitorSyncCmd.Dir = frontendPath
	capacitorSyncCmd.Stdout = os.Stdout
	capacitorSyncCmd.Stderr = os.Stderr

	if err := capacitorSyncCmd.Run(); err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao sincronizar Capacitor: %v\n", err)
	}

	// 4. Build do iOS usando Xcode (requer Xcode)
	fmt.Println("🍎 Para compilar o app iOS, abra o projeto no Xcode:")
	fmt.Printf("   open %s/App/App.xcworkspace\n", iosPath)
	fmt.Println("\nOu use o comando:")
	fmt.Printf("   cd %s/App && xcodebuild -workspace App.xcworkspace -scheme App -configuration Release -archivePath App.xcarchive archive\n", iosPath)

	fmt.Printf("\n✓ Build concluído!\n")
	fmt.Printf("📱 Projeto iOS: %s\n", iosPath)

	return nil
}

func buildDesktop(projectConfig *config.ProjectConfig) error {
	fmt.Println("🔨 Compilando projeto Desktop...")

	// Verificar se frontend existe
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	// 1. Build do servidor Go
	fmt.Println("📦 Compilando servidor Go...")

	// Configurar CGO baseado no banco de dados
	goBuildCmd := exec.Command("go", "build", "-o", "server", "cmd/server/main.go")
	if projectConfig.Database == "sqlite" {
		// Para SQLite, desabilitar CGO para usar modernc.org/sqlite (puro Go)
		goBuildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	}
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr

	if err := goBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar servidor Go: %w", err)
	}

	// 2. Copiar binário para pasta do Electron
	electronPath := filepath.Join(frontendPath, "src-electron")
	if err := os.MkdirAll(electronPath, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório electron: %w", err)
	}

	// Copiar binário para electron
	serverDest := filepath.Join(electronPath, "server")
	if runtime.GOOS == "windows" {
		serverDest = filepath.Join(electronPath, "server.exe")
	}

	if err := copyFile("server", serverDest); err != nil {
		return fmt.Errorf("erro ao copiar binário para Electron: %w", err)
	}

	// Se usar SQLite, copiar banco de dados para src-electron (se existir)
	if projectConfig.Database == "sqlite" {
		dbName := projectConfig.ProjectName
		if dbName == "" {
			dbName = "app"
		}
		dbFile := dbName + ".db"

		// Procurar arquivo .db no diretório atual ou em subdiretórios comuns
		var dbPath string
		possiblePaths := []string{
			dbFile,
			filepath.Join("data", dbFile),
			filepath.Join("database", dbFile),
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				dbPath = path
				break
			}
		}

		if dbPath != "" {
			fmt.Printf("📦 Copiando banco SQLite embutido: %s\n", dbPath)
			dbDest := filepath.Join(electronPath, "database.db")
			if err := copyFile(dbPath, dbDest); err != nil {
				fmt.Printf("⚠️  Aviso: Erro ao copiar banco SQLite: %v\n", err)
			} else {
				fmt.Println("✓ Banco SQLite copiado para Electron")
			}
		} else {
			fmt.Println("ℹ️  Nenhum banco SQLite encontrado - será criado em tempo de execução")
		}
	}

	// Copiar arquivo .env para src-electron/ se existir (para o extraResources do electron-builder)
	envSource := ".env"
	if _, err := os.Stat(envSource); err == nil {
		envDest := filepath.Join(electronPath, ".env")
		if err := copyFile(envSource, envDest); err != nil {
			fmt.Printf("⚠️  Aviso: Erro ao copiar .env para src-electron/: %v\n", err)
		} else {
			fmt.Println("✓ Arquivo .env copiado para src-electron/")
		}
	} else {
		fmt.Println("ℹ️  Arquivo .env não encontrado na raiz do projeto")
	}

	// 3. Build do Quasar Electron
	fmt.Println("📦 Compilando Quasar Electron...")
	quasarBuildCmd := exec.Command("quasar", "build", "-m", "electron")
	quasarBuildCmd.Dir = frontendPath
	quasarBuildCmd.Stdout = os.Stdout
	quasarBuildCmd.Stderr = os.Stderr

	if err := quasarBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar Quasar Electron: %w", err)
	}

	// 4. Com electron-builder, os arquivos são copiados automaticamente via extraResources
	// Mas vamos verificar se foram copiados corretamente e informar ao usuário
	// O electron-builder gera estrutura diferente: dist/electron/win-unpacked/ ou dist/electron/

	distPath := filepath.Join(frontendPath, "dist", "electron")

	// Verificar estrutura do builder (win-unpacked)
	winUnpackedPath := filepath.Join(distPath, "win-unpacked")
	if _, err := os.ReadDir(winUnpackedPath); err == nil {
		// Estrutura do electron-builder encontrada
		resourcesPath := filepath.Join(winUnpackedPath, "resources")

		// Verificar se os arquivos foram copiados pelo extraResources
		serverExePath := filepath.Join(resourcesPath, "server.exe")
		if runtime.GOOS != "windows" {
			serverExePath = filepath.Join(resourcesPath, "server")
		}

		if _, err := os.Stat(serverExePath); err == nil {
			fmt.Printf("✓ Servidor Go encontrado em resources/ (copiado pelo electron-builder)\n")
		} else {
			// Fallback: copiar manualmente se não foi copiado pelo builder
			fmt.Println("⚠️  Servidor Go não encontrado em resources/, copiando manualmente...")
			if err := os.MkdirAll(resourcesPath, 0755); err == nil {
				serverSource := filepath.Join(electronPath, "server.exe")
				if runtime.GOOS != "windows" {
					serverSource = filepath.Join(electronPath, "server")
				}
				if _, err := os.Stat(serverSource); err == nil {
					serverDest := filepath.Join(resourcesPath, filepath.Base(serverSource))
					if err := copyFile(serverSource, serverDest); err == nil {
						fmt.Printf("✓ Servidor Go copiado manualmente para resources/\n")
					}
				}
			}
		}

		// Verificar banco e .env
		if projectConfig.Database == "sqlite" {
			dbPath := filepath.Join(resourcesPath, "database.db")
			if _, err := os.Stat(dbPath); err != nil {
				// Copiar manualmente se não foi copiado
				dbSource := filepath.Join(electronPath, "database.db")
				if _, err := os.Stat(dbSource); err == nil {
					if err := copyFile(dbSource, dbPath); err == nil {
						fmt.Printf("✓ Banco SQLite copiado manualmente para resources/\n")
					}
				}
			}
		}

		envPath := filepath.Join(resourcesPath, ".env")
		if _, err := os.Stat(envPath); err != nil {
			// Copiar manualmente se não foi copiado
			envSource := ".env"
			if _, err := os.Stat(envSource); err == nil {
				if err := copyFile(envSource, envPath); err == nil {
					fmt.Printf("✓ Arquivo .env copiado manualmente para resources/\n")
				}
			}
		}
	} else {
		// Estrutura antiga do packager (Packaged)
		packagedPath := filepath.Join(distPath, "Packaged")
		if entries, err := os.ReadDir(packagedPath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					resourcesPath := filepath.Join(packagedPath, entry.Name(), "resources")
					if err := os.MkdirAll(resourcesPath, 0755); err == nil {
						// Copiar servidor Go
						serverSource := filepath.Join(electronPath, "server.exe")
						if runtime.GOOS != "windows" {
							serverSource = filepath.Join(electronPath, "server")
						}
						if _, err := os.Stat(serverSource); err == nil {
							serverDest := filepath.Join(resourcesPath, filepath.Base(serverSource))
							if err := copyFile(serverSource, serverDest); err == nil {
								fmt.Printf("✓ Servidor Go copiado para resources/: %s\n", serverDest)
							}
						}

						// Copiar banco SQLite se existir
						if projectConfig.Database == "sqlite" {
							dbSource := filepath.Join(electronPath, "database.db")
							if _, err := os.Stat(dbSource); err == nil {
								dbDest := filepath.Join(resourcesPath, "database.db")
								if err := copyFile(dbSource, dbDest); err == nil {
									fmt.Printf("✓ Banco SQLite copiado para resources/: %s\n", dbDest)
								}
							}
						}

						// Copiar arquivo .env se existir
						envSource := ".env"
						if _, err := os.Stat(envSource); err == nil {
							envDest := filepath.Join(resourcesPath, ".env")
							if err := copyFile(envSource, envDest); err == nil {
								fmt.Printf("✓ Arquivo .env copiado para resources/: %s\n", envDest)
							} else {
								fmt.Printf("⚠️  Aviso: Erro ao copiar .env: %v\n", err)
							}
						} else {
							fmt.Println("ℹ️  Arquivo .env não encontrado - usando variáveis de ambiente do sistema")
						}
					}
				}
			}
		}
	}

	// Limpar binário temporário
	os.Remove("server")

	// O instalador .exe será gerado pelo Electron Builder
	// distPath já foi declarado acima

	fmt.Printf("✓ Build concluído!\n")
	fmt.Printf("📁 Diretório dist: %s\n", distPath)

	// Procurar pelo instalador .exe gerado pelo electron-builder (NSIS)
	installerFiles, _ := filepath.Glob(filepath.Join(distPath, "*.exe"))
	if len(installerFiles) > 0 {
		fmt.Printf("💾 Instalador gerado: %s\n", installerFiles[0])
	} else {
		// Procurar também por executável não empacotado
		exeFiles, _ := filepath.Glob(filepath.Join(distPath, "**", "*.exe"))
		if len(exeFiles) > 0 {
			fmt.Printf("💾 Executável: %s\n", exeFiles[0])
		}
	}

	fmt.Println("ℹ️  Com electron-builder, os arquivos (server.exe, database.db, .env) são incluídos automaticamente via extraResources")

	return nil
}

func buildWeb(projectConfig *config.ProjectConfig) error {
	fmt.Println("🔨 Compilando projeto Web...")

	// Criar diretório build
	buildDir := "build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório build: %w", err)
	}

	// 1. Build do servidor Go
	fmt.Println("📦 Compilando servidor Go...")
	goBuildCmd := exec.Command("go", "build", "-o", filepath.Join(buildDir, "server"), "cmd/server/main.go")
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr

	if err := goBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar servidor Go: %w", err)
	}

	// 2. Build do Quasar SPA
	fmt.Println("📦 Compilando Quasar SPA...")
	frontendPath := projectConfig.FrontendDir
	if frontendPath == "" {
		frontendPath = "frontend"
	}

	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		return fmt.Errorf("diretório frontend não encontrado")
	}

	quasarBuildCmd := exec.Command("quasar", "build")
	quasarBuildCmd.Dir = frontendPath
	quasarBuildCmd.Stdout = os.Stdout
	quasarBuildCmd.Stderr = os.Stderr

	if err := quasarBuildCmd.Run(); err != nil {
		return fmt.Errorf("erro ao compilar Quasar: %w", err)
	}

	// 3. Copiar dist/spa para build/spa
	spaSource := filepath.Join(frontendPath, "dist", "spa")
	spaDest := filepath.Join(buildDir, "spa")

	fmt.Println("📋 Copiando arquivos SPA para build...")
	if err := copyDir(spaSource, spaDest); err != nil {
		return fmt.Errorf("erro ao copiar arquivos SPA: %w", err)
	}

	fmt.Printf("✓ Build concluído!\n")
	fmt.Printf("📁 Diretório build: %s\n", buildDir)
	fmt.Printf("  - Binário Go: %s/server\n", buildDir)
	fmt.Printf("  - SPA: %s/spa\n", buildDir)

	return nil
}

// copyDir copia um diretório recursivamente
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copia um arquivo
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
	if err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// createMobilePackage cria um package Go básico para gomobile
func createMobilePackage() error {
	mobileDir := "cmd/mobile"
	if err := os.MkdirAll(mobileDir, 0755); err != nil {
		return err
	}

	content := `package mobile

// Export é uma função exemplo exportada para Android
func Export() string {
	return "Gaver Mobile"
}
`

	return os.WriteFile(filepath.Join(mobileDir, "mobile.go"), []byte(content), 0644)
}

// configureGradleForAAR configura o build.gradle para incluir o AAR
func configureGradleForAAR(buildGradlePath, libsPath string) error {
	// Ler o arquivo build.gradle
	content, err := os.ReadFile(buildGradlePath)
	if err != nil {
		return fmt.Errorf("erro ao ler build.gradle: %w", err)
	}

	contentStr := string(content)

	// Verificar se já tem a configuração do flatDir
	hasFlatDir := false
	if contains(contentStr, "flatDir") {
		hasFlatDir = true
	}

	// Verificar se já tem a dependência do AAR
	hasAARDependency := false
	if contains(contentStr, "gaver.aar") || contains(contentStr, "libs/gaver") {
		hasAARDependency = true
	}

	// Se já está configurado, não precisa fazer nada
	if hasFlatDir && hasAARDependency {
		return nil
	}

	// Adicionar repositório flatDir se não existir
	if !hasFlatDir {
		// Procurar pela seção repositories
		if contains(contentStr, "repositories {") {
			// Adicionar flatDir dentro de repositories existente
			contentStr = addFlatDirToRepositories(contentStr)
		} else {
			// Adicionar seção repositories completa
			// Procurar onde inserir (geralmente após android {)
			if contains(contentStr, "android {") {
				contentStr = insertAfter(contentStr, "android {", "\n    repositories {\n        flatDir {\n            dirs 'libs'\n        }\n    }\n")
			}
		}
	}

	// Adicionar dependência do AAR se não existir
	if !hasAARDependency {
		// Procurar pela seção dependencies
		if contains(contentStr, "dependencies {") {
			// Adicionar dependência dentro de dependencies existente
			contentStr = addAARDependency(contentStr)
		} else {
			// Adicionar seção dependencies completa
			// Procurar onde inserir (geralmente antes do final do arquivo ou após android {})
			contentStr = insertBeforeClosingBrace(contentStr, "    dependencies {\n        implementation(name: 'gaver', ext: 'aar')\n    }\n")
		}
	}

	// Escrever o arquivo modificado
	return os.WriteFile(buildGradlePath, []byte(contentStr), 0644)
}

// Funções auxiliares para manipulação de strings
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsInMiddle(s, substr)))
}

func containsInMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func addFlatDirToRepositories(content string) string {
	// Procurar por "repositories {" e adicionar flatDir dentro
	reposIndex := findStringIndex(content, "repositories {")
	if reposIndex == -1 {
		return content
	}

	// Encontrar o final do bloco repositories
	braceCount := 0
	insertIndex := reposIndex + len("repositories {")
	for i := insertIndex; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				// Inserir flatDir antes do fechamento
				before := content[:i]
				after := content[i:]
				return before + "        flatDir {\n            dirs 'libs'\n        }\n    " + after
			}
		}
	}

	return content
}

func addAARDependency(content string) string {
	// Procurar por "dependencies {" e adicionar dependência dentro
	depsIndex := findStringIndex(content, "dependencies {")
	if depsIndex == -1 {
		return content
	}

	// Encontrar o final do bloco dependencies
	braceCount := 0
	insertIndex := depsIndex + len("dependencies {")
	for i := insertIndex; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
		} else if content[i] == '}' {
			braceCount--
			if braceCount == 0 {
				// Inserir dependência antes do fechamento
				before := content[:i]
				after := content[i:]
				return before + "    implementation(name: 'gaver', ext: 'aar')\n    " + after
			}
		}
	}

	return content
}

func insertAfter(content, search, insert string) string {
	index := findStringIndex(content, search)
	if index == -1 {
		return content
	}
	insertPos := index + len(search)
	return content[:insertPos] + insert + content[insertPos:]
}

func insertBeforeClosingBrace(content, insert string) string {
	// Encontrar o último "}" que fecha o bloco principal
	lastBraceIndex := -1
	braceCount := 0
	for i := len(content) - 1; i >= 0; i-- {
		if content[i] == '}' {
			braceCount++
			if lastBraceIndex == -1 {
				lastBraceIndex = i
			}
		} else if content[i] == '{' {
			braceCount--
			if braceCount == 0 && lastBraceIndex != -1 {
				return content[:lastBraceIndex] + insert + "}" + content[lastBraceIndex+1:]
			}
		}
	}
	return content
}

func findStringIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// generateAndroidKeystore gera um keystore para assinar o APK Android
func generateAndroidKeystore(androidPath string, projectConfig *config.ProjectConfig) error {
	keystorePath := filepath.Join(androidPath, "app", "keystore.jks")
	keystorePropsPath := filepath.Join(androidPath, "app", ".keystore.properties")

	// Verificar se keystore já existe
	if _, err := os.Stat(keystorePath); err == nil {
		fmt.Println("✓ Keystore já existe, usando o existente")
		return nil
	}

	fmt.Println("🔑 Gerando keystore para assinatura Android...")

	// Verificar se keytool está disponível
	if _, err := exec.LookPath("keytool"); err != nil {
		return fmt.Errorf("keytool não encontrado. Instale o JDK para gerar keystore")
	}

	// Valores padrão para o keystore
	keystoreAlias := "key"
	keystorePassword := "gaver123"
	keyPassword := "gaver123"
	validity := "10000" // ~27 anos
	commonName := projectConfig.ProjectName
	if commonName == "" {
		commonName = "Gaver App"
	}

	// Gerar keystore
	keytoolCmd := exec.Command("keytool",
		"-genkeypair",
		"-v",
		"-storetype", "PKCS12",
		"-keystore", keystorePath,
		"-alias", keystoreAlias,
		"-keyalg", "RSA",
		"-keysize", "2048",
		"-validity", validity,
		"-storepass", keystorePassword,
		"-keypass", keyPassword,
		"-dname", fmt.Sprintf("CN=%s, OU=Development, O=Gaver, L=Unknown, ST=Unknown, C=BR", commonName),
	)
	keytoolCmd.Stdout = os.Stdout
	keytoolCmd.Stderr = os.Stderr

	if err := keytoolCmd.Run(); err != nil {
		return fmt.Errorf("erro ao gerar keystore: %w", err)
	}

	fmt.Println("✓ Keystore gerado com sucesso")

	// Salvar propriedades do keystore em arquivo
	keystoreProps := fmt.Sprintf(`storeFile=keystore.jks
storePassword=%s
keyAlias=%s
keyPassword=%s
`, keystorePassword, keystoreAlias, keyPassword)

	if err := os.WriteFile(keystorePropsPath, []byte(keystoreProps), 0600); err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao salvar .keystore.properties: %v\n", err)
		fmt.Println("   Você precisará configurar manualmente o build.gradle")
	} else {
		fmt.Println("✓ Propriedades do keystore salvas em .keystore.properties")
	}

	// Configurar build.gradle para usar o keystore
	buildGradlePath := filepath.Join(androidPath, "app", "build.gradle")
	if err := configureGradleForSigning(buildGradlePath, keystorePropsPath); err != nil {
		fmt.Printf("⚠️  Aviso: Erro ao configurar signing no build.gradle: %v\n", err)
		fmt.Println("   Você precisará configurar manualmente o signing no build.gradle")
	} else {
		fmt.Println("✓ build.gradle configurado para usar keystore")
	}

	return nil
}

// configureGradleForSigning configura o build.gradle para usar o keystore
func configureGradleForSigning(buildGradlePath, keystorePropsPath string) error {
	content, err := os.ReadFile(buildGradlePath)
	if err != nil {
		return fmt.Errorf("erro ao ler build.gradle: %w", err)
	}

	contentStr := string(content)

	// Verificar se já tem signing config
	if strings.Contains(contentStr, "signingConfigs") {
		return nil // Já configurado
	}

	// Ler propriedades do keystore
	propsContent, err := os.ReadFile(keystorePropsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler .keystore.properties: %w", err)
	}

	// Parsear propriedades
	props := make(map[string]string)
	for _, line := range strings.Split(string(propsContent), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			props[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	storeFile := props["storeFile"]
	storePassword := props["storePassword"]
	keyAlias := props["keyAlias"]
	keyPassword := props["keyPassword"]

	if storeFile == "" || storePassword == "" || keyAlias == "" || keyPassword == "" {
		return fmt.Errorf("propriedades do keystore incompletas")
	}

	// Adicionar signing config antes do bloco android
	signingConfig := fmt.Sprintf(`
    signingConfigs {
        release {
            storeFile file('%s')
            storePassword '%s'
            keyAlias '%s'
            keyPassword '%s'
        }
    }
`, storeFile, storePassword, keyAlias, keyPassword)

	// Procurar bloco android e adicionar signingConfigs
	if idx := findStringIndex(contentStr, "android {"); idx != -1 {
		insertPos := idx + len("android {")
		contentStr = contentStr[:insertPos] + signingConfig + contentStr[insertPos:]
	} else {
		return fmt.Errorf("não foi possível encontrar bloco 'android {' no build.gradle")
	}

	// Adicionar signingConfig ao buildTypes release
	releaseConfig := `
        release {
            signingConfig signingConfigs.release
        }
`
	if strings.Contains(contentStr, "buildTypes {") {
		// Procurar buildTypes e adicionar release
		buildTypesIdx := findStringIndex(contentStr, "buildTypes {")
		if buildTypesIdx != -1 {
			insertPos := buildTypesIdx + len("buildTypes {")
			// Verificar se já tem release
			if !strings.Contains(contentStr[insertPos:], "release {") {
				contentStr = contentStr[:insertPos] + releaseConfig + contentStr[insertPos:]
			}
		}
	}

	return os.WriteFile(buildGradlePath, []byte(contentStr), 0644)
}

// generateMainActivity gera o MainActivity.java que inicia o servidor Go
func generateMainActivity(androidPath string, projectConfig *config.ProjectConfig) error {
	// Determinar o package name (geralmente do capacitor.config.js ou do build.gradle)
	// Por padrão, usar um nome baseado no nome do projeto
	packageName := "com.gaver.app"

	// Tentar ler do capacitor.config.js
	capacitorConfigPath := filepath.Join(filepath.Dir(androidPath), "capacitor.config.js")
	if content, err := os.ReadFile(capacitorConfigPath); err == nil {
		// Procurar por appId no capacitor.config.js
		contentStr := string(content)
		if idx := findStringIndex(contentStr, "appId:"); idx != -1 {
			// Extrair o appId
			start := idx + len("appId:")
			end := start
			for end < len(contentStr) && (contentStr[end] == ' ' || contentStr[end] == '\'' || contentStr[end] == '"') {
				end++
			}
			start = end
			for end < len(contentStr) && contentStr[end] != '\'' && contentStr[end] != '"' && contentStr[end] != ',' && contentStr[end] != '\n' {
				end++
			}
			if end > start {
				packageName = contentStr[start:end]
			}
		}
	}

	// Ler template usando o sistema de templates embarcado
	templateContent, err := templates.TemplatesFS.ReadFile("main_activity_android.tmpl")
	if err != nil {
		return fmt.Errorf("erro ao ler template MainActivity: %w", err)
	}

	// Parsear e executar template
	tmpl, err := template.New("main_activity").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("erro ao parsear template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"PackageName": packageName,
		"ProjectName": projectConfig.ProjectName,
	}); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}
	content := buf.String()

	// Determinar caminho do MainActivity
	// Normalmente está em android/app/src/main/java/com/.../MainActivity.java
	// Mas o package pode variar, então vamos procurar ou criar em um local padrão
	mainActivityDir := filepath.Join(androidPath, "app", "src", "main", "java", strings.ReplaceAll(packageName, ".", string(filepath.Separator)))
	if err := os.MkdirAll(mainActivityDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório MainActivity: %w", err)
	}

	mainActivityPath := filepath.Join(mainActivityDir, "MainActivity.java")
	return os.WriteFile(mainActivityPath, []byte(content), 0644)
}
