package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/modules"

	"github.com/spf13/cobra"
)

func NewModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module",
		Short: "Gerencia módulos do projeto",
		Long:  "Cria e gerencia módulos com estrutura completa (models, handlers, services, repositories).",
	}

	// Subcomandos
	cmd.AddCommand(newModuleCreateCommand())
	cmd.AddCommand(newModuleModelCommand())
	cmd.AddCommand(newModuleCrudCommand())

	return cmd
}

func newModuleCreateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create [nome]",
		Short: "Cria um novo módulo",
		Long:  "Cria a estrutura completa de um módulo com pastas para models, handlers, services e repositories.",
		Args:  cobra.ExactArgs(1),
		RunE:  runModuleCreate,
	}
}

func runModuleCreate(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	fmt.Printf("Criando módulo '%s'...\n", moduleName)

	if err := modules.CreateModule(moduleName); err != nil {
		return fmt.Errorf("erro ao criar módulo: %w", err)
	}

	fmt.Printf("✓ Módulo '%s' criado com sucesso!\n\n", moduleName)
	fmt.Println("Estrutura criada:")
	fmt.Printf("  modules/%s/\n", moduleName)
	fmt.Println("  ├── models/")
	fmt.Println("  ├── handlers/")
	fmt.Println("  ├── services/")
	fmt.Println("  ├── repositories/")
	fmt.Println("  ├── validators/")
	fmt.Println("  └── module.go")
	fmt.Println("\nPróximos passos:")
	fmt.Printf("  gaver module:model %s User name:string email:string\n", moduleName)
	fmt.Printf("  gaver module:crud %s User\n", moduleName)

	return nil
}

func newModuleModelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "model [module] [ModelName] [fields...]",
		Short: "Cria um model dentro de um módulo",
		Long:  "Gera um arquivo de model com annotations gaverModel.",
		Example: `  gaver module model users User name:string email:string:unique age:int
  gaver module model products Product title:string price:float stock:int`,
		Args: cobra.MinimumNArgs(2),
		RunE: runModuleModel,
	}
}

func runModuleModel(cmd *cobra.Command, args []string) error {
	moduleName := args[0]
	modelName := args[1]
	fields := []string{}
	if len(args) > 2 {
		fields = args[2:]
	}

	fmt.Printf("Gerando model '%s' no módulo '%s'...\n", modelName, moduleName)

	if err := modules.CreateModel(moduleName, modelName, fields); err != nil {
		return fmt.Errorf("erro ao criar model: %w", err)
	}

	fmt.Printf("✓ Model '%s' criado em modules/%s/models/%s.go\n", modelName, moduleName, toLower(modelName))

	return nil
}

func newModuleCrudCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crud [module] [ModelName]",
		Short: "Gera CRUD completo para um model",
		Long:  "Gera handlers, services e repositories com operações CRUD completas.",
		Args:  cobra.ExactArgs(2),
		RunE:  runModuleCrud,
	}

	cmd.Flags().StringSlice("only", []string{}, "Gera apenas os métodos especificados (list,get,create,update,delete)")
	cmd.Flags().StringSlice("except", []string{}, "Gera todos exceto os métodos especificados")

	return cmd
}

func runModuleCrud(cmd *cobra.Command, args []string) error {
	moduleName := args[0]
	modelName := args[1]

	only, _ := cmd.Flags().GetStringSlice("only")
	except, _ := cmd.Flags().GetStringSlice("except")

	fmt.Printf("Gerando CRUD para '%s' no módulo '%s'...\n", modelName, moduleName)

	if err := modules.CreateCRUD(moduleName, modelName, only, except); err != nil {
		return fmt.Errorf("erro ao criar CRUD: %w", err)
	}

	fmt.Printf("✓ CRUD gerado com sucesso!\n\n")
	fmt.Println("Arquivos criados:")
	fmt.Printf("  - modules/%s/handlers/%s_handler.go\n", moduleName, toLower(modelName))
	fmt.Printf("  - modules/%s/services/%s_service.go\n", moduleName, toLower(modelName))
	fmt.Printf("  - modules/%s/repositories/%s_repository.go\n", moduleName, toLower(modelName))

	return nil
}

func toLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]+32) + s[1:]
}
