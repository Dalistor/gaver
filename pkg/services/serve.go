package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/Dalistor/gaver/pkg/types"
)

func Serve(gaverModuleFile *types.GaverModuleFile, port string) error {
	switch gaverModuleFile.Type {
	case "api":
		return serveAPI(gaverModuleFile, port)
	}

	return nil
}

func serveAPI(gaverModuleFile *types.GaverModuleFile, port string) error {
	if port == "" {
		port = "7077"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar handler para capturar Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nEncerrando servidor...")
		cancel()
	}()

	fmt.Println("Iniciando servidor API na porta", port)
	fmt.Println("Servidor disponível em http://localhost:" + port)
	fmt.Println("Pressione Ctrl+C para parar o servidor")

	cmd := exec.CommandContext(ctx, "go", "run", "cmd/api/main.go", "-p", port)

	if gaverModuleFile.ProjectDatabaseType == "sqlite" {
		if err := checkCompilerInstalled(); err != nil {
			return fmt.Errorf("erro ao verificar compilador C: %w", err)
		}

		cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.Canceled {
			fmt.Println("Servidor encerrado pelo usuário")
			return nil
		}
		return fmt.Errorf("erro ao executar servidor: %w", err)
	}

	return nil
}

func checkCompilerInstalled() error {
	compilers := []string{"gcc", "clang", "cc"}
	compilerFound := false

	for _, compiler := range compilers {
		if _, err := exec.LookPath(compiler); err == nil {
			compilerFound = true
			break
		}
	}

	if !compilerFound {
		return fmt.Errorf("compilador C não encontrado. Instale um compilador C (gcc, clang ou cc) para usar SQLite com CGO.")
	}

	return nil
}
