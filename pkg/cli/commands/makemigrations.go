package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/services"
	"github.com/spf13/cobra"
)

func NewMakeMigrationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "makemigrations",
		Short: "Gera as migrações para o banco de dados",
		Long:  `Gera as migrações para o banco de dados`,
		RunE:  makemigrationsCommand,
	}

	return cmd
}

func makemigrationsCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	gaverModuleFile, err := services.ReadGaverModuleFile()
	if err != nil {
		return fmt.Errorf("erro ao obter o arquivo module: %w", err)
	}

	if err := services.MakeMigrations(gaverModuleFile); err != nil {
		return fmt.Errorf("erro ao gerar as migrações: %w", err)
	}

	return nil
}
