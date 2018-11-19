package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
)

var commands []subcommands.Command

func main() {
	flag.Parse()
	subcommands.Register(subcommands.HelpCommand(), "")
	for _, c := range commands {
		subcommands.Register(c, "")
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}
