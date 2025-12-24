package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/services"
	"github.com/Dalistor/gaver/pkg/types"
	"github.com/spf13/cobra"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init -n [nome do projeto] -d [banco de dados] -p [tipo de projeto]",
		Short: "Inicia um novo projeto Gaver",
		Long:  `Inicia um novo projeto Gaver com uma estrutura de pastas e arquivos padronizados`,
		RunE:  initCommand,
	}

	cmd.Flags().StringP("name", "n", "", "Nome do projeto")
	cmd.Flags().StringP("database", "d", "", "Banco de dados a ser usado (mysql, postgres ou sqlite)")
	cmd.Flags().StringP("project-type", "p", "", "Tipo de projeto a ser criado (web, api, mobile ou desktop)")
	cmd.Flags().BoolP("help", "h", false, "Ajuda com o comando init")

	return cmd
}

func initCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	initCommand := &types.InitCommand{}
	initCommand.Name = cmd.Flag("name").Value.String()
	initCommand.Database = cmd.Flag("database").Value.String()
	initCommand.ProjectType = cmd.Flag("project-type").Value.String()

	if err := validateInitCommand(initCommand); err != nil {
		return fmt.Errorf("erro ao validar o comando init: %w. \n\nUse --help para mais informações.", err)
	}

	// fazer download do template
	fmt.Println("Fazendo download do template")
	if err := services.Download(initCommand, nil); err != nil {
		return fmt.Errorf("erro ao fazer download do template: %w", err)
	}
	fmt.Println("Template baixado com sucesso")

	// parsear arquivos
	fmt.Println("Parseando arquivos")
	if err := services.Parse(initCommand); err != nil {
		return fmt.Errorf("erro ao parsear arquivos: %w", err)
	}
	fmt.Println("Arquivos parseados com sucesso")

	fmt.Println("Base do projeto criada com sucesso")
	fmt.Println("Para iniciar o projeto, use o comando: 'gaver serve'")
	fmt.Println("Tenha um bom desenvolvimento!")

	return nil
}

func validateInitCommand(initCommand *types.InitCommand) error {
	if initCommand.Name == "" {
		return fmt.Errorf("nome do projeto é obrigatório")
	}

	if initCommand.Database == "" {
		return fmt.Errorf("banco de dados é obrigatório")
	}

	if initCommand.ProjectType == "" {
		return fmt.Errorf("tipo de projeto é obrigatório")
	}

	if initCommand.ProjectType != "web" && initCommand.ProjectType != "api" && initCommand.ProjectType != "mobile" && initCommand.ProjectType != "desktop" {
		return fmt.Errorf("tipo de projeto inválido")
	}

	if initCommand.Database != "mysql" && initCommand.Database != "postgres" && initCommand.Database != "sqlite" {
		return fmt.Errorf("banco de dados inválido")
	}

	if initCommand.ProjectType == "web" {
		if initCommand.Database != "mysql" && initCommand.Database != "postgres" {
			return fmt.Errorf("banco de dados mysql ou postgres é obrigatório para projeto web")
		}
	}

	fmt.Println("Iniciando um novo projeto Gaver")
	fmt.Println("Nome do projeto:", initCommand.Name)
	fmt.Println("Banco de dados:", initCommand.Database)
	fmt.Println("Tipo de projeto:", initCommand.ProjectType)

	return nil
}
