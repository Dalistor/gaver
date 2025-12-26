package commands

import (
	"fmt"
	"strings"

	"github.com/Dalistor/gaver/pkg/services"
	"github.com/Dalistor/gaver/pkg/types"
	"github.com/spf13/cobra"
)

var TEMPLATE_FILES = map[string]string{
	"repository": "https://raw.githubusercontent.com/Dalistor/Gaver-Modules/refs/heads/module/repositories/repository.tmplt_crud",
	"service":    "https://raw.githubusercontent.com/Dalistor/Gaver-Modules/refs/heads/module/services/service.tmplt_crud",
	"handler":    "https://raw.githubusercontent.com/Dalistor/Gaver-Modules/refs/heads/module/handlers/handler.tmplt_crud",
	"route":      "https://raw.githubusercontent.com/Dalistor/Gaver-Modules/refs/heads/module/routes/route.tmplt_crud",
}

func NewModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module -n [nome do modulo] -t [tipo de modulo] -c [Gerar Repository, Service, Routes baseado no model]",
		Short: "Gerencia os módulos do projeto",
		Long:  `Gerencia os módulos do projeto com uma interface de linha de comando`,
		RunE:  moduleCommand,
	}

	cmd.Flags().StringP("name", "n", "", "Nome do modulo")
	cmd.Flags().StringP("type", "t", "", "Tipo de modulo [crud, service=name]")
	cmd.Flags().StringP("controller", "c", "", "Gerar Controller [nome do model]")

	return cmd
}

func moduleCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	moduleCommand := &types.ModuleCommand{
		Name:       cmd.Flag("name").Value.String(),
		Type:       cmd.Flag("type").Value.String(),
		Controller: cmd.Flag("controller").Value.String(),
	}

	if err := validateModuleCommand(moduleCommand); err != nil {
		return fmt.Errorf("erro ao validar o comando module: %w. \n\nUse --help para mais informações.", err)
	}

	// Ler arquivo module
	gaverModuleFile, err := services.ReadGaverModuleFile()
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo module: %w", err)
	}

	if moduleCommand.Controller == "" {
		// Adicionar modulo ao arquivo module
		gaverModuleFile.ProjectModules = append(gaverModuleFile.ProjectModules, moduleCommand.Name)
		if err := services.SetGaverModuleFile(gaverModuleFile); err != nil {
			return fmt.Errorf("erro ao adicionar modulo ao arquivo module: %w", err)
		}

		// fazer download do template
		fmt.Println("Fazendo download do template")
		if err := services.DownloadTemplate("module", moduleCommand.Name); err != nil {
			return fmt.Errorf("erro ao fazer download do template: %w", err)
		}
		fmt.Println("Template baixado com sucesso")

		// parsear arquivos
		fmt.Println("Parseando arquivos")
		if err := services.ParseModule(moduleCommand, gaverModuleFile.ProjectName); err != nil {
			return fmt.Errorf("erro ao parsear arquivos: %w", err)
		}
		fmt.Println("Arquivos parseados com sucesso")
	} else {
		// Fazer downlaod dos templates de CRUD
		fmt.Println("Fazendo download dos templates de CRUD")

		modelNameInLowercase := strings.ToLower(moduleCommand.Controller)

		if err := services.DownloadSingleFileFromGit(TEMPLATE_FILES["repository"], fmt.Sprintf("modules/%s/repositories", moduleCommand.Name), modelNameInLowercase+".tmplt_crud"); err != nil {
			return fmt.Errorf("erro ao fazer download do template de repository: %w", err)
		}
		if err := services.DownloadSingleFileFromGit(TEMPLATE_FILES["service"], fmt.Sprintf("modules/%s/services", moduleCommand.Name), modelNameInLowercase+".tmplt_crud"); err != nil {
			return fmt.Errorf("erro ao fazer download do template de service: %w", err)
		}
		if err := services.DownloadSingleFileFromGit(TEMPLATE_FILES["handler"], fmt.Sprintf("modules/%s/handlers", moduleCommand.Name), modelNameInLowercase+".tmplt_crud"); err != nil {
			return fmt.Errorf("erro ao fazer download do template de handler: %w", err)
		}

		fmt.Println("Templates de CRUD baixados com sucesso")

		// parsear arquivos
		fmt.Println("Parseando arquivos")
		if err := services.ParseModule(moduleCommand, gaverModuleFile.ProjectName); err != nil {
			return fmt.Errorf("erro ao parsear arquivos: %w", err)
		}
		fmt.Println("Arquivos parseados com sucesso")

	}

	return nil
}

func validateModuleCommand(moduleCommand *types.ModuleCommand) error {
	if moduleCommand.Name == "" {
		return fmt.Errorf("nome do modulo é obrigatório")
	}

	if moduleCommand.Controller != "" && moduleCommand.Name == "" {
		return fmt.Errorf("nome do modulo é obrigatório para gerar controller")
	}

	if moduleCommand.Controller != "" && moduleCommand.Type != "" {
		return fmt.Errorf("tipo de modulo não pode ser informado para gerar controller")
	}

	return nil
}
