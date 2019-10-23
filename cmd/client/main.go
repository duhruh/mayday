package main

import (
	"fmt"
	"os"

	"github.com/docker/mayday/internal/cmd"
)

func main() {
	command := cmd.NewRootCommand()
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
