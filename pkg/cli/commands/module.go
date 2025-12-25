package commands

import (
	"fmt"
	"strings"

	"github.com/Dalistor/gaver/pkg/services"
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
	cmd.Flags().StringP("type", "t", "", "Tipo de modulo [crud, service=name]")
	cmd.Flags().BoolP("controller", "c", false, "Gerar Controller [isto gerará repository, service, route e handler baseado no model]")

	return cmd
}

func moduleCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	moduleCommand := &types.ModuleCommand{
		Name:       cmd.Flag("name").Value.String(),
		Type:       cmd.Flag("type").Value.String(),
		Controller: cmd.Flag("controller").Changed,
	}

	if err := validateModuleCommand(moduleCommand); err != nil {
		return fmt.Errorf("erro ao validar o comando module: %w. \n\nUse --help para mais informações.", err)
	}

	// Ler arquivo module
	gaverModuleFile, err := services.ReadGaverModuleFile()
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo module: %w", err)
	}

	// Adicionar modulo ao arquivo module
	gaverModuleFile.ProjectModules = append(gaverModuleFile.ProjectModules, moduleCommand.Name)
	if err := services.SetGaverModuleFile(gaverModuleFile); err != nil {
		return fmt.Errorf("erro ao adicionar modulo ao arquivo module: %w", err)
	}

	// fazer download do template
	fmt.Println("Fazendo download do template")
	if err := services.DownloadAPI("module", moduleCommand.Name); err != nil {
		return fmt.Errorf("erro ao fazer download do template: %w", err)
	}
	fmt.Println("Template baixado com sucesso")

	return nil
}

func validateModuleCommand(moduleCommand *types.ModuleCommand) error {
	if moduleCommand.Name == "" {
		return fmt.Errorf("nome do modulo é obrigatório")
	}

	if !strings.Contains(moduleCommand.Type, "crud") && !strings.Contains(moduleCommand.Type, "service") {
		return fmt.Errorf("tipo de modulo inválido")
	}

	return nil
}
