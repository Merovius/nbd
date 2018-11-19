package main

import (
	"context"
	"flag"
	"os"
	"strconv"

	"github.com/google/subcommands"
)

var commands []subcommands.Command

func main() {
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		subcommands.ImportantFlag(f.Name)
	})
	subcommands.Register(subcommands.HelpCommand(), "")
	for _, c := range commands {
		subcommands.Register(c, "")
	}
	os.Exit(int(subcommands.Execute(context.Background())))
}

type indexFlag struct {
	set bool
	val uint32
	def string
}

func (f *indexFlag) String() string {
	if f.set {
		return strconv.FormatUint(uint64(f.val), 10)
	}
	if f.def != "" {
		return f.def
	}
	return "auto"
}

func (f *indexFlag) Set(s string) error {
	def := f.def
	if def == "" {
		def = "auto"
	}
	if s == def {
		f.set = false
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	f.set = true
	f.val = uint32(v)
	return nil
}
