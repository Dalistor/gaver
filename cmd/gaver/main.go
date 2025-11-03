package main

import (
	"fmt"
	"os"

	"github.com/Dalistor/gaver/pkg/cli"
)

func main() {
	app := cli.NewCLI()

	if err := app.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro: %v\n", err)
		os.Exit(1)
	}
}
