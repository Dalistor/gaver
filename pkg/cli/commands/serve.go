package commands

import (
	"fmt"

	"github.com/Dalistor/gaver/pkg/services"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve -p [porta]",
		Short: "Inicia o servidor de desenvolvimento",
		Long:  `Inicia o servidor de desenvolvimento para o projeto Gaver`,
		RunE:  serveCommand,
	}

	cmd.Flags().StringP("port", "p", "7077", "Porta do servidor")
	cmd.Flags().BoolP("help", "h", false, "Ajuda com o comando serve")

	return cmd
}

func serveCommand(cmd *cobra.Command, args []string) error {
	if cmd.Flag("help").Changed {
		return cmd.Help()
	}

	gaverModuleFile, err := services.ReadGaverModuleFile()
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo gaverModule.json: %w\n\nCertifique-se de que o .env esteja configurado corretamente.", err)
	}

	port := cmd.Flag("port").Value.String()
	if err := services.Serve(gaverModuleFile, port); err != nil {
		return fmt.Errorf("erro ao iniciar o servidor: %w", err)
	}

	return nil
}
