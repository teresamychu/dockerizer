// Dockerize CLI - generates Dockerfiles for Go Codebases
package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/teresamychu/dockerizer/commands"
)

func main() {
	c := cli.NewCLI("dockerize", "0.1.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"generate": func() (cli.Command, error) {
			return &commands.GenerateCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(exitStatus)
}
