package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Dalistor/gaver/pkg/config"
	"github.com/Dalistor/gaver/pkg/generator/structure"

	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [nome-do-projeto] -d [tipo-de-banco] -t [tipo-de-projeto]",
		Short: "Inicializa um novo projeto Gaver",
		Long:  "Cria um novo projeto Gaver com a estrutura de diretórios e arquivos padrão.\nTipos de projeto: server (padrão), android, desktop, web",
		Args:  cobra.ExactArgs(1),
		RunE:  run_init,
	}

	// Adicionar flags
	cmd.Flags().StringP("database", "d", "mysql", "Tipo de banco (postgres, mysql, sqlite)")
	cmd.Flags().StringP("type", "t", "server", "Tipo de projeto (server, android, desktop, web)")

	return cmd
}

func run_init(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	database, _ := cmd.Flags().GetString("database")
	projectType, _ := cmd.Flags().GetString("type")

	// Validar tipo de projeto
	if !config.IsValidProjectType(projectType) {
		return fmt.Errorf("tipo de projeto inválido: %s. Use: server, android, desktop ou web", projectType)
	}

	fmt.Printf("Inicializando projeto: %s (tipo: %s)...\n", projectName, projectType)

	// Criar configuração do projeto
	projectConfig := &config.ProjectConfig{
		ProjectName: projectName,
		Type:        config.ProjectType(projectType),
		Database:    database,
		ServerPort:  "8080",
		FrontendDir: "frontend",
	}

	// Gerar arquivos base
	if err := structure.GenerateInitialFiles(projectName, database, projectType); err != nil {
		return fmt.Errorf("erro ao gerar arquivos: %w", err)
	}

	// Gerar arquivos específicos do tipo
	switch config.ProjectType(projectType) {
	case config.ProjectTypeAndroid:
		if err := setupAndroidProject(projectName, projectConfig); err != nil {
			return fmt.Errorf("erro ao configurar projeto Android: %w", err)
		}
	case config.ProjectTypeDesktop:
		if err := setupDesktopProject(projectName, projectConfig); err != nil {
			return fmt.Errorf("erro ao configurar projeto Desktop: %w", err)
		}
	case config.ProjectTypeWeb:
		if err := setupWebProject(projectName, projectConfig); err != nil {
			return fmt.Errorf("erro ao configurar projeto Web: %w", err)
		}
	}

	// Escrever configuração do projeto
	if err := config.WriteProjectConfig(projectConfig, projectName); err != nil {
		return fmt.Errorf("erro ao escrever configuração: %w", err)
	}

	fmt.Println("✓ Arquivos iniciais gerados")

	fmt.Printf("\n✓ Projeto '%s' inicializado com sucesso!\n\n", projectName)
	fmt.Println("Próximos passos:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")

	// Nota: npm install já foi executado automaticamente para projetos com frontend
	if projectType == "server" {
		fmt.Println("\nPara começar a desenvolver:")
		fmt.Println("  gaver serve")
	} else {
		fmt.Println("\nPara começar a desenvolver:")
		fmt.Println("  gaver serve")
	}

	return nil
}

func setupAndroidProject(projectName string, projectConfig *config.ProjectConfig) error {
	fmt.Println("Configurando projeto Android...")

	// Verificar se gomobile está instalado
	if _, err := exec.LookPath("gomobile"); err != nil {
		fmt.Println("Instalando gomobile...")
		installCmd := exec.Command("go", "install", "golang.org/x/mobile/cmd/gomobile@latest")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("erro ao instalar gomobile: %w", err)
		}
		fmt.Println("✓ gomobile instalado")
	}

	// Gerar estrutura frontend Android
	if err := structure.GenerateAndroidFrontend(projectName, projectConfig); err != nil {
		return fmt.Errorf("erro ao gerar frontend Android: %w", err)
	}

	return nil
}

func setupDesktopProject(projectName string, projectConfig *config.ProjectConfig) error {
	fmt.Println("Configurando projeto Desktop...")

	// Gerar estrutura frontend Desktop
	if err := structure.GenerateDesktopFrontend(projectName, projectConfig); err != nil {
		return fmt.Errorf("erro ao gerar frontend Desktop: %w", err)
	}

	return nil
}

func setupWebProject(projectName string, projectConfig *config.ProjectConfig) error {
	fmt.Println("Configurando projeto Web...")

	// Gerar estrutura frontend Web
	if err := structure.GenerateWebFrontend(projectName, projectConfig); err != nil {
		return fmt.Errorf("erro ao gerar frontend Web: %w", err)
	}

	return nil
}
