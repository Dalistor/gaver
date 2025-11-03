package cli

import (
	"github.com/Dalistor/gaver/pkg/cli/commands"

	"github.com/spf13/cobra"
)

type CLI struct {
	RootCmd *cobra.Command
}

// Nova instancia do CLI
func NewCLI() *CLI {
	cli := &CLI{}

	cli.RootCmd = &cobra.Command{
		Use:   "gaver",
		Short: "Gaver - Framework web em golang",
		Long:  "Gaver é um framework web em golang que facilita a criação de aplicativos web.",
	}

	cli.registerCommands()

	return cli
}

// Executa o CLI
func (cli *CLI) Execute() error {
	return cli.RootCmd.Execute()
}

// Registra comandos
func (cli *CLI) registerCommands() {
	cli.RootCmd.AddCommand(commands.NewInitCommand())
	cli.RootCmd.AddCommand(commands.NewModuleCommand())
	cli.RootCmd.AddCommand(commands.NewMigrationsCommand())
	cli.RootCmd.AddCommand(commands.NewMigrateCommand())
	cli.RootCmd.AddCommand(commands.NewServeCommand())
}
