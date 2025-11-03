package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/generator/structure"

	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [nome-do-projeto] -d [tipo-de-banco]",
		Short: "Inicializa um novo projeto Gaver",
		Long:  "Cria um novo projeto Gaver com a estrutura de diretórios e arquivos padrão.",
		Args:  cobra.ExactArgs(1),
		RunE:  run_init,
	}

	// Adicionar flags
	cmd.Flags().StringP("database", "d", "mysql", "Tipo de banco (postgres, mysql, sqlite)")

	return cmd
}

func run_init(cmd *cobra.Command, args []string) error {
	projectName := args[0]
	database, _ := cmd.Flags().GetString("database")

	fmt.Printf("Inicializando projeto: %s...\n", projectName)

	// Gerar arquivos
	if err := structure.GenerateInitialFiles(projectName, database); err != nil {
		return fmt.Errorf("erro ao gerar arquivos: %w", err)
	}

	fmt.Println("✓ Arquivos iniciais gerados")

	fmt.Printf("\n✓ Projeto '%s' inicializado com sucesso!\n\n", projectName)
	fmt.Println("Próximos passos:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  go mod tidy")
	fmt.Println("  gaver module create users")
	fmt.Println("  gaver module model users User name:string email:string")
	fmt.Println("  gaver module crud users User")
	fmt.Println("  gaver serve")

	return nil
}
