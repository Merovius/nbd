package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/subcommands"
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	subcommands.Register(new(listCmd), "")
	subcommands.Register(new(serveCmd), "")
	os.Exit(int(subcommands.Execute(context.Background())))
}
