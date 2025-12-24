package commands

import (
	"fmt"
	"strings"

	"github.com/Dalistor/gaver/pkg/types"
	"github.com/spf13/cobra"
)

func NewModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module -n [nome do modulo] -t [tipo de modulo] -c [Gerar Controller, Repository, Service, Routes baseado no model]",
		Short: "Gerencia os módulos do projeto",
		Long:  `Gerencia os módulos do projeto com uma interface de linha de comando`,
		RunE:  moduleCommand,
	}

	cmd.Flags().StringP("name", "n", "", "Nome do modulo")
	cmd.Flags().StringP("type", "t", "", "Tipo de modulo")
	cmd.Flags().BoolP("controller", "c", false, "Gerar Controller")

	return cmd
}

func moduleCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	moduleCommand := &types.ModuleCommand{
		Name: cmd.Flag("name").Value.String(),
		Type: cmd.Flag("type").Value.String(),
		Controller: types.List[string]{
			Items: strings.Split(cmd.Flag("controller").Value.String(), ","),
		},
	}

	if err := validateModuleCommand(moduleCommand); err != nil {
		return fmt.Errorf("erro ao validar o comando module: %w. \n\nUse --help para mais informações.", err)
	}

	return nil
}

func validateModuleCommand(moduleCommand *types.ModuleCommand) error {
	if moduleCommand.Name == "" {
		return fmt.Errorf("nome do modulo é obrigatório")
	}

	if !strings.Contains(moduleCommand.Type, "crud") && !strings.Contains(moduleCommand.Type, "service") {
		return fmt.Errorf("tipo de modulo inválido")
	}

	for _, controller := range moduleCommand.Controller.Items {
		if controller != "repository" && controller != "service" && controller != "route" && controller != "handler" {
			return fmt.Errorf("controller %s inválido", controller)
		}
	}

	return nil
}
