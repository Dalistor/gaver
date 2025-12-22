package main

import (
	"github.com/Dalistor/gaver/pkg/cli"

	"fmt"
	"os"
)

func main() {
	cli := cli.NewCLI()
	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
