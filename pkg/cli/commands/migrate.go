package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/services"
	"github.com/spf13/cobra"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Executa as migrações para o banco de dados",
		Long:  `Executa as migrações para o banco de dados`,
		RunE:  migrateCommand,
	}

	return cmd
}

func migrateCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	gaverModuleFile, err := services.ReadGaverModuleFile()
	if err != nil {
		return fmt.Errorf("erro ao obter o arquivo module: %w", err)
	}

	if err := services.Migrate(gaverModuleFile); err != nil {
		return fmt.Errorf("erro ao executar as migrações: %w", err)
	}

	return nil
}
