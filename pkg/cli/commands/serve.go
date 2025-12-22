package commands

import (
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve -p [porta]",
		Short: "Inicia o servidor de desenvolvimento",
		Long:  `Inicia o servidor de desenvolvimento para o projeto Gaver`,
		RunE:  serveCommand,
	}

	cmd.Flags().StringP("port", "p", "8080", "Porta do servidor")
	cmd.Flags().BoolP("help", "h", false, "Ajuda com o comando serve")

	return cmd
}

func serveCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	return nil
}
