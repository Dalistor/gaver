package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Dalistor/gaver/pkg/config"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Inicia o servidor de desenvolvimento",
		Long:  "Executa o servidor da aplicação em modo de desenvolvimento. Para projetos Android/Desktop, também inicia o Quasar dev server.",
		RunE:  runServe,
	}

	cmd.Flags().Bool("android", false, "Abre Android Studio para debug (apenas para projetos Android)")

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
		return runNormal()
	}

	// Executar baseado no tipo de projeto
	switch projectConfig.Type {
	case config.ProjectTypeAndroid:
		openAndroidStudio, _ := cmd.Flags().GetBool("android")
		return runAndroid(projectConfig, openAndroidStudio)
	case config.ProjectTypeDesktop:
		return runDesktop(projectConfig)
	case config.ProjectTypeWeb:
		return runWeb(projectConfig)
	case config.ProjectTypeServer:
		return runNormal()
	default:
		return runNormal()
	}
}

func runNormal() error {
	fmt.Printf("Iniciando servidor...\n")
	fmt.Println("Use Ctrl+C para parar o servidor")

	// Executar go run
	cmd := exec.Command("go", "run", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

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

func runAndroid(projectConfig *config.ProjectConfig, openAndroidStudio bool) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

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

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Iniciar Quasar dev server
	quasarCmd := exec.Command("quasar", "dev", "-m", "capacitor", "-T", "android")
	quasarCmd.Dir = frontendPath
	quasarCmd.Stdout = os.Stdout
	quasarCmd.Stderr = os.Stderr

	if err := quasarCmd.Start(); err != nil {
		goServer.Process.Kill()
		return fmt.Errorf("erro ao iniciar Quasar dev server: %w", err)
	}

	// Se flag -android, abrir Android Studio
	if openAndroidStudio {
		androidPath := filepath.Join(frontendPath, "android")
		if _, err := os.Stat(androidPath); err == nil {
			fmt.Println("📱 Abrindo Android Studio...")
			studioCmd := exec.Command("studio", androidPath)
			studioCmd.Start() // Não esperar, apenas iniciar
		}
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

func runDesktop(projectConfig *config.ProjectConfig) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

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

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Iniciar Quasar dev server
	quasarCmd := exec.Command("quasar", "dev", "-m", "electron")
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

func runWeb(projectConfig *config.ProjectConfig) error {
	fmt.Println("🚀 Iniciando servidor Go e Quasar dev server...")
	fmt.Println("Use Ctrl+C para parar os servidores")

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

	if err := goServer.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar servidor Go: %w", err)
	}

	// Iniciar Quasar dev server
	quasarCmd := exec.Command("quasar", "dev")
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
