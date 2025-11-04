package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Inicia o servidor de desenvolvimento",
		Long:  "Executa o servidor da aplicação em modo de desenvolvimento.",
		RunE:  runServe,
	}

	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	// Verificar se estamos em um projeto Gaver
	if _, err := os.Stat("cmd/server/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("não parece ser um projeto Gaver. Execute 'gaver init' primeiro")
	}

	return runNormal()
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
