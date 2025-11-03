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

	cmd.Flags().StringP("port", "p", "8080", "Porta do servidor")
	cmd.Flags().BoolP("watch", "w", false, "Watch mode (hot reload)")

	return cmd
}

func runServe(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")
	watch, _ := cmd.Flags().GetBool("watch")

	// Verificar se estamos em um projeto Gaver
	if _, err := os.Stat("cmd/server/main.go"); os.IsNotExist(err) {
		return fmt.Errorf("não parece ser um projeto Gaver. Execute 'gaver init' primeiro")
	}

	// Setar variável de ambiente para a porta
	os.Setenv("SERVER_PORT", port)

	if watch {
		return runWithWatch(port)
	}

	return runNormal(port)
}

func runNormal(port string) error {
	fmt.Printf("🚀 Iniciando servidor na porta %s...\n", port)
	fmt.Println("📝 Use Ctrl+C para parar o servidor")
	fmt.Println()

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

func runWithWatch(port string) error {
	// TODO: Implementar hot reload
	fmt.Println("⚠️  Watch mode ainda não implementado")
	fmt.Println("    Usando modo normal...")
	return runNormal(port)
}

