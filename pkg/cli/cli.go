package cli

import (
	"github.com/Dalistor/gaver/pkg/cli/commands"
	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd *cobra.Command
}

func NewCLI() *CLI {
	cli := &CLI{
		rootCmd: &cobra.Command{
			Use:   "gaver",
			Short: "Gaver - Framework multiplataformas",
			Long:  `Gaver é um framework multiplataformas para criação de aplicações web, mobile e desktop`,
		},
	}

	cli.registerCommands(cli.rootCmd)
	return cli
}

func (c *CLI) Execute() error {
	return c.rootCmd.Execute()
}

func (c *CLI) registerCommands(cmd *cobra.Command) {
	cmd.AddCommand(commands.NewInitCommand())
	cmd.AddCommand(commands.NewServeCommand())
}
