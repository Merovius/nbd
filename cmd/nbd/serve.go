package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Merovius/nbd"
	"github.com/google/subcommands"
)

type serveCmd struct {
	addr string
	unix bool
}

func (cmd *serveCmd) Name() string {
	return "serve"
}

func (cmd *serveCmd) Synopsis() string {
	return "serve NBD devices and their status"
}

func (cmd *serveCmd) Usage() string {
	return `Usage: nbd serve <file>

Serve a file as a network block device.
`
}

func (cmd *serveCmd) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&cmd.addr, "addr", "localhost:10809", "Address to listen on")
	fs.BoolVar(&cmd.unix, "unix", false, "Treat -addr as a unix domain socket")
}

func (cmd *serveCmd) Execute(ctx context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 1 {
		log.Print(cmd.Usage())
		return subcommands.ExitUsageError
	}

	f, err := os.OpenFile(fs.Arg(0), os.O_RDWR, 0)
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	log.Println(fi.Size())

	idx, wait, err := nbd.Loopback(ctx, f, uint64(fi.Size()))
	if err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	fmt.Printf("Connected to /dev/nbd%d\n", idx)
	if err := wait(); err != nil {
		log.Println(err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
